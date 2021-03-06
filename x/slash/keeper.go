package slash

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/crypto"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey        sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc             *codec.Codec // The wire codec for binary encoding/decoding.
	validatorKeeper validator.Keeper
	paramstore      params.Subspace
}

// NewKeeper creates new instances of the slash Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, validatorKeeper validator.Keeper, paramstore params.Subspace) Keeper {
	return Keeper{
		storeKey:        storeKey,
		cdc:             cdc,
		validatorKeeper: validatorKeeper,
		paramstore:      paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// HandleGuardFailure handles a validator fails to guard state.
func (k Keeper) HandleGuardFailure(ctx sdk.Context, beneficiaryAddr, failedAddr sdk.AccAddress) {
	failedValAddr := sdk.ValAddress(failedAddr)
	failedValidator, found := k.validatorKeeper.GetValidator(ctx, failedValAddr)
	if !found {
		log.Errorf("Cannot find failed validator %s", failedValAddr)
		return
	}

	var beneficiaries []AccountFractionPair
	// TODO: need to add address(0) as the miningPool and make sure the total share is 1
	if !beneficiaryAddr.Empty() {
		beneficiaryValAddr := sdk.ValAddress(beneficiaryAddr)
		beneficiaryValidator, found := k.validatorKeeper.GetValidator(ctx, beneficiaryValAddr)
		if !found {
			log.Errorf("Cannot find beneficiary validator %s", beneficiaryValAddr)
			return
		}
		beneficiaries = append(beneficiaries, NewAccountFractionPair(beneficiaryValidator.Description.Identity, k.SlashFractionGuardFailure(ctx)))
	}

	k.Slash(ctx, AttributeValueGuardFailure, failedValidator, calculateSlashAmount(failedValidator.GetConsensusPower(), k.SlashFractionGuardFailure(ctx)), beneficiaries)
}

// HandleProposalDepositBurn handles a depositor supports refused proposal.
func (k Keeper) HandleProposalDepositBurn(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Int) {
	valAddr := sdk.ValAddress(depositor)
	validator, found := k.validatorKeeper.GetValidator(ctx, valAddr)
	if !found {
		log.Errorf("Cannot find failed validator %s", valAddr)
		return
	}

	k.Slash(ctx, AttributeValueDepositBurn, validator, amount.MulRaw(common.TokenDec), []AccountFractionPair{})
}

// HandleDoubleSign handles a validator signing two blocks at the same height.
// power: power of the double-signing validator at the height of infraction
func (k Keeper) HandleDoubleSign(ctx sdk.Context, addr crypto.Address, power int64) {
	consAddr := sdk.ConsAddress(addr)
	validator, found := k.validatorKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		log.Errorf("Cannot find validator %s", consAddr)
		return
	}

	log.Infof("Confirmed double sign from %s", consAddr)
	k.Slash(ctx, slashing.AttributeValueDoubleSign, validator, calculateSlashAmount(power, k.SlashFractionDoubleSign(ctx)), []AccountFractionPair{})
}

// HandleValidatorSignature handles a validator signature, must be called once per validator per block.
func (k Keeper) HandleValidatorSignature(ctx sdk.Context, addr crypto.Address, power int64, signed bool) {
	height := ctx.BlockHeight()
	consAddr := sdk.ConsAddress(addr)
	validator, found := k.validatorKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		log.Errorf("Cannot find validator %s", consAddr)
		return
	}

	signInfo, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		signInfo = slashing.NewValidatorSigningInfo(
			consAddr,
			height,
			0,
			time.Unix(0, 0),
			false,
			0,
		)
	}

	// this is a relative index, so it counts blocks the validator *should* have signed
	// will use the 0-value default signing info if not present, except for start height
	signedBlocksWindow := k.SignedBlocksWindow(ctx)
	index := signInfo.IndexOffset % signedBlocksWindow
	signInfo.IndexOffset++

	// Update signed block bit array & counter
	// This counter just tracks the sum of the bit array
	// That way we avoid needing to read/write the whole array each time
	previous := k.GetValidatorMissedBlockBitArray(ctx, consAddr, index)
	missed := !signed
	switch {
	case !previous && missed:
		// Array value has changed from not missed to missed, increment counter
		k.SetValidatorMissedBlockBitArray(ctx, consAddr, index, true)
		signInfo.MissedBlocksCounter++
	case previous && !missed:
		// Array value has changed from missed to not missed, decrement counter
		k.SetValidatorMissedBlockBitArray(ctx, consAddr, index, false)
		signInfo.MissedBlocksCounter--
	default:
		// Array value at this index has not changed, no need to update counter
	}

	minHeight := signInfo.StartHeight + signedBlocksWindow
	maxMissed := signedBlocksWindow - k.MinSignedPerWindow(ctx).MulInt64(signedBlocksWindow).RoundInt64()

	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {
		// Downtime confirmed: slash the validator
		log.Infof("Validator %s past min height of %d and above max miss threshold of %d",
			consAddr, minHeight, maxMissed)

		// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
		signInfo.MissedBlocksCounter = 0
		signInfo.IndexOffset = 0
		k.ClearValidatorMissedBlockBitArray(ctx, consAddr)
		k.Slash(ctx, slashing.AttributeValueMissingSignature, validator, calculateSlashAmount(power, k.SlashFractionDowntime(ctx)), []AccountFractionPair{})
	}

	k.SetValidatorSigningInfo(ctx, signInfo)
}

// Slash a validator for an infraction
// Find the contributing stake and burn the specified slashFactor of it
func (k Keeper) Slash(ctx sdk.Context, reason string, failedValidator staking.Validator, slashAmount sdk.Int, beneficiaries []AccountFractionPair) {
	candidate, found := k.validatorKeeper.GetCandidate(ctx, failedValidator.Description.Identity)
	if !found {
		log.Errorln("Cannot find candidate profile for the failed validator", failedValidator.Description.Identity)
		return
	}

	penalty := NewPenalty(k.GetNextPenaltyNonce(ctx), reason, failedValidator.Description.Identity)
	for _, delegator := range candidate.Delegators {
		penaltyAmt := slashAmount.Mul(delegator.DelegatedStake).Quo(candidate.StakingPool)
		accountAmtPair := NewAccountAmtPair(delegator.DelegatorAddr, penaltyAmt)
		penalty.PenalizedDelegators = append(penalty.PenalizedDelegators, accountAmtPair)
	}

	penalty.Beneficiaries = beneficiaries
	penalty.GenerateProtoBytes()
	k.SetPenalty(ctx, penalty)

	log.Warnf("Slash validator: %s, amount: %s, reason: %s, nonce: %d",
		failedValidator.GetOperator(), slashAmount, reason, penalty.Nonce)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeSlash,
			sdk.NewAttribute(sdk.AttributeKeyAction, ActionPenalty),
			sdk.NewAttribute(AttributeKeyNonce, sdk.NewUint(penalty.Nonce).String()),
			sdk.NewAttribute(slashing.AttributeKeyReason, reason),
		),
	)
}

// Gets the next Penalty nonce, and increment nonce by 1
func (k Keeper) GetNextPenaltyNonce(ctx sdk.Context) (nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(PenaltyNonceKey)
	if bz != nil {
		nonce = binary.BigEndian.Uint64(bz)
	}

	store.Set(PenaltyNonceKey, sdk.Uint64ToBigEndian(nonce+1))
	return
}

// Gets the entire Penalty metadata for a nonce
func (k Keeper) GetPenalty(ctx sdk.Context, nonce uint64) (penalty Penalty, found bool) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(GetPenaltyKey(nonce)) {
		return penalty, false
	}

	value := store.Get(GetPenaltyKey(nonce))
	k.cdc.MustUnmarshalBinaryBare(value, &penalty)
	return penalty, true
}

// Sets the entire Penalty metadata for a nonce
func (k Keeper) SetPenalty(ctx sdk.Context, penalty Penalty) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetPenaltyKey(penalty.Nonce), k.cdc.MustMarshalBinaryBare(penalty))
}

// Stored by *validator* address (not operator address)
func (k Keeper) GetValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress) (info slashing.ValidatorSigningInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(slashing.GetValidatorSigningInfoKey(address))
	if bz == nil {
		found = false
		return
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
	found = true
	return
}

// Stored by *validator* address (not operator address)
func (k Keeper) SetValidatorSigningInfo(ctx sdk.Context, info slashing.ValidatorSigningInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(info)
	store.Set(slashing.GetValidatorSigningInfoKey(info.Address), bz)
}

// Stored by *validator* address (not operator address)
func (k Keeper) GetValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, index int64) (missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(slashing.GetValidatorMissedBlockBitArrayKey(address, index))
	if bz == nil {
		// lazy: treat empty key as not missed
		missed = false
		return
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &missed)
	return
}

// Stored by *validator* address (not operator address)
func (k Keeper) SetValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, index int64, missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(missed)
	store.Set(slashing.GetValidatorMissedBlockBitArrayKey(address, index), bz)
}

// Stored by *validator* address (not operator address)
func (k Keeper) ClearValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, slashing.GetValidatorMissedBlockBitArrayPrefixKey(address))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

func calculateSlashAmount(power int64, slashFactor sdk.Dec) sdk.Int {
	if slashFactor.IsNegative() {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}
	// Amount of slashing = slash slashFactor * power at time of infraction
	amount := sdk.TokensFromConsensusPower(power).MulRaw(common.TokenDec)
	return amount.ToDec().Mul(slashFactor).TruncateInt()
}
