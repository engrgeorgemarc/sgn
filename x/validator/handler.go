package validator

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// NewHandler returns a handler for "validator" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgClaimValidator:
			return handleMsgClaimValidator(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized validator Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set eth address
func handleMsgClaimValidator(ctx sdk.Context, keeper Keeper, msg MsgClaimValidator) sdk.Result {
	cp, err := keeper.ethClient.Guard.CandidateProfiles(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(msg.EthAddress))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query candidate profile: %s", err)).Result()
	}

	if !sdk.AccAddress(cp.SidechainAddr).Equals(msg.Sender) {
		return sdk.ErrInternal("Sender is not selected validator").Result()
	}

	valAddress := sdk.ValAddress(msg.Sender)
	validator, found := keeper.stakingKeeper.GetValidator(ctx, valAddress)
	_, f := keeper.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(msg.PubKey))

	if found != f {
		return sdk.ErrInternal("Invalid sender address or public key").Result()
	}

	if !found {
		description := staking.Description{
			Moniker: msg.EthAddress,
		}
		validator = staking.NewValidator(valAddress, msg.PubKey, description)
	}

	validator, _ = validator.AddTokensFromDel(sdk.NewIntFromBigInt(cp.Stakes))
	keeper.stakingKeeper.SetValidator(ctx, validator)
	keeper.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	keeper.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)

	return sdk.Result{}
}
