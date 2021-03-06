package validator

import (
	"fmt"

	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "validator" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		logEntry := seal.NewMsgLog()
		var res *sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgSetTransactors:
			res, err = handleMsgSetTransactors(ctx, keeper, msg, logEntry)
		case MsgWithdrawReward:
			res, err = handleMsgWithdrawReward(ctx, keeper, msg, logEntry)
		case MsgSignReward:
			res, err = handleMsgSignReward(ctx, keeper, msg, logEntry)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}

		if err != nil {
			logEntry.Error = append(logEntry.Error, err.Error())
		}

		seal.CommitMsgLog(logEntry)
		return res, err
	}
}

// Handle a message to set transactors
func handleMsgSetTransactors(ctx sdk.Context, keeper Keeper, msg MsgSetTransactors, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	for _, transactor := range msg.Transactors {
		logEntry.Transactor = append(logEntry.Transactor, transactor.String())
	}

	candidate, found := keeper.GetCandidate(ctx, msg.EthAddress)
	if !found {
		return nil, fmt.Errorf("Candidate does not exist")
	}

	if !candidate.Operator.Equals(msg.Sender) {
		return nil, fmt.Errorf("The candidate is not operated by the sender.")
	}

	candidate.Transactors = msg.Transactors
	for _, transactor := range candidate.Transactors {
		keeper.InitAccount(ctx, transactor)
	}

	keeper.SetCandidate(ctx, candidate)
	return &sdk.Result{}, nil
}

// Handle a message to withdraw reward
func handleMsgWithdrawReward(ctx sdk.Context, keeper Keeper, msg MsgWithdrawReward, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return nil, fmt.Errorf("Reward does not exist")
	}
	if !reward.HasNewReward() {
		return nil, fmt.Errorf("No new reward")
	}

	reward.InitateWithdraw()
	keeper.SetReward(ctx, reward)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			ModuleName,
			sdk.NewAttribute(sdk.AttributeKeyAction, ActionInitiateWithdraw),
			sdk.NewAttribute(AttributeKeyEthAddress, msg.EthAddress),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// Handle a message to sign reward
func handleMsgSignReward(ctx sdk.Context, keeper Keeper, msg MsgSignReward, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	validator, found := keeper.stakingKeeper.GetValidator(ctx, sdk.ValAddress(msg.Sender))
	if !found {
		return nil, fmt.Errorf("Sender is not validator")
	}
	if validator.Status != sdk.Bonded {
		return nil, fmt.Errorf("Validator is not bonded")
	}

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return nil, fmt.Errorf("Reward does not exist")
	}

	err := reward.AddSig(msg.Sig, validator.Description.Identity)
	if err != nil {
		return nil, fmt.Errorf("Failed to add sig: %s", err)
	}

	keeper.SetReward(ctx, reward)
	return &sdk.Result{}, nil
}
