package keeper

import (
	"fmt"

	"github.com/celer-network/sgn/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetDeposit gets the deposit of a specific depositor on a specific proposal
func (keeper Keeper) GetDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) (deposit types.Deposit, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.DepositKey(proposalID, depositorAddr))
	if bz == nil {
		return deposit, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
	return deposit, true
}

// SetDeposit sets a Deposit to the gov store
func (keeper Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(types.DepositKey(deposit.ProposalID, deposit.Depositor), bz)
}

// GetAllDeposits returns all the deposits from the store
func (keeper Keeper) GetAllDeposits(ctx sdk.Context) (deposits types.Deposits) {
	keeper.IterateAllDeposits(ctx, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// GetDeposits returns all the deposits from a proposal
func (keeper Keeper) GetDeposits(ctx sdk.Context, proposalID uint64) (deposits types.Deposits) {
	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// IterateAllDeposits iterates over the all the stored deposits and performs a callback function
func (keeper Keeper) IterateAllDeposits(ctx sdk.Context, cb func(deposit types.Deposit) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DepositsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// IterateDeposits iterates over the all the proposals deposits and performs a callback function
func (keeper Keeper) IterateDeposits(ctx sdk.Context, proposalID uint64, cb func(deposit types.Deposit) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DepositsKey(proposalID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// GetDepositor gets the depositor by address
func (keeper Keeper) GetDepositor(ctx sdk.Context, depositorAddr sdk.AccAddress) (depositor types.Depositor, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.DepositorKey(depositorAddr))
	if bz == nil {
		return depositor, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &depositor)
	return depositor, true
}

// SetDepositor sets a Depositor to the gov store
func (keeper Keeper) SetDepositor(ctx sdk.Context, depositorAddr sdk.AccAddress, depositor types.Depositor) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(depositor)
	store.Set(types.DepositorKey(depositorAddr), bz)
}

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (keeper Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Int) (bool, error) {
	// Checks to see if proposal exists
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return false, sdkerrors.Wrapf(types.ErrUnknownProposal, "%d", proposalID)
	}

	// Check if proposal is still depositable
	if (proposal.Status != types.StatusDepositPeriod) && (proposal.Status != types.StatusVotingPeriod) {
		return false, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "%d", proposalID)
	}

	validator, found := keeper.vk.GetValidator(ctx, sdk.ValAddress(depositorAddr))
	if !found {
		return false, sdkerrors.Wrapf(types.ErrUnknownProposal, "Invalid depositor addr %s", depositorAddr)
	}

	depositor, found := keeper.GetDepositor(ctx, depositorAddr)
	if !found {
		depositor = types.NewDepositor()
	}

	if ctx.BlockTime().Before(depositor.MutedUntil) {
		return false, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Depositor is muted")
	}

	depositor.Amount = depositor.Amount.Add(depositAmount)
	if depositor.Amount.GT(validator.Tokens) {
		return false, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Depositor does not have enough stake to deposit")
	}
	keeper.SetDepositor(ctx, depositorAddr, depositor)

	// Update proposal
	proposal.TotalDeposit = proposal.TotalDeposit.Add(depositAmount)
	keeper.SetProposal(ctx, proposal)

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	if proposal.Status == types.StatusDepositPeriod && proposal.TotalDeposit.GTE(keeper.GetDepositParams(ctx).MinDeposit) {
		keeper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	// Add or update deposit object
	deposit, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
	if found {
		deposit.Amount = deposit.Amount.Add(depositAmount)
	} else {
		deposit = types.NewDeposit(proposalID, depositorAddr, depositAmount)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalDeposit,
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	keeper.SetDeposit(ctx, deposit)
	return activatedVotingPeriod, nil
}

// RefundDeposits refunds and deletes all the deposits on a specific proposal
func (keeper Keeper) RefundDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) bool {
		depositor, _ := keeper.GetDepositor(ctx, deposit.Depositor)
		depositor.Amount = depositor.Amount.Sub(deposit.Amount)
		keeper.SetDepositor(ctx, deposit.Depositor, depositor)

		store.Delete(types.DepositKey(proposalID, deposit.Depositor))
		return false
	})
}

// DeleteDeposits deletes all the deposits on a specific proposal without refunding them
func (keeper Keeper) DeleteDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) bool {
		// TODO: properly handle delete deposits
		depositor, _ := keeper.GetDepositor(ctx, deposit.Depositor)
		depositor.Amount = depositor.Amount.Sub(deposit.Amount)
		depositor.MutedUntil = ctx.BlockTime().Add(keeper.GetDepositParams(ctx).MutedDuration)
		keeper.SetDepositor(ctx, deposit.Depositor, depositor)
		keeper.sk.HandleProposalDepositBurn(ctx, deposit.Depositor, deposit.Amount)

		store.Delete(types.DepositKey(proposalID, deposit.Depositor))
		return false
	})
}
