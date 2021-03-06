package cli

import (
	"fmt"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/x/validator/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingCli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

const (
	flagSeq = "seq"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	validatorQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the validator module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	validatorQueryCmd.AddCommand(flags.GetCommands(
		GetCmdPuller(storeKey, cdc),
		GetCmdDelegator(storeKey, cdc),
		GetCmdCandidate(storeKey, cdc),
		GetCmdReward(storeKey, cdc),
		GetCmdRewardRequest(storeKey, cdc),
		stakingCli.GetCmdQueryValidator(staking.StoreKey, cdc),
		stakingCli.GetCmdQueryValidators(staking.StoreKey, cdc),
	)...)
	return validatorQueryCmd
}

// GetCmdPuller queries puller info
func GetCmdPuller(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "puller",
		Short: "query puller info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			puller, err := QueryPuller(cliCtx, queryRoute)
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(puller)
		},
	}
}

// Query puller info
func QueryPuller(cliCtx context.CLIContext, queryRoute string) (puller types.Puller, err error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPuller)
	res, err := common.RobustQuery(cliCtx, route)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &puller)
	return
}

// GetCmdPusher queries pusher info
func GetCmdPusher(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pusher",
		Short: "query pusher info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			pusher, err := QueryPusher(cliCtx, queryRoute)
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(pusher)
		},
	}
}

// Query pusher info
func QueryPusher(cliCtx context.CLIContext, queryRoute string) (pusher types.Pusher, err error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPusher)
	res, err := common.RobustQuery(cliCtx, route)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &pusher)
	return
}

// GetCmdDelegator queries request info
func GetCmdDelegator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delegator [candidateAddress] [delegatorAddress]",
		Short: "query delegator info by candidateAddress and delegatorAddress",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			delegator, err := QueryDelegator(cliCtx, queryRoute, args[0], args[1])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(delegator)
		},
	}
}

func QueryDelegator(cliCtx context.CLIContext, queryRoute, candidateAddress, delegatorAddress string) (delegator types.Delegator, err error) {
	data, err := cliCtx.Codec.MarshalJSON(types.NewQueryDelegatorParams(candidateAddress, delegatorAddress))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegator)
	res, err := common.RobustQueryWithData(cliCtx, route, data)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &delegator)
	return
}

// GetCmdCandidate queries request info
func GetCmdCandidate(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "candidate [candidateAddress]",
		Short: "query candidate info by candidateAddress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			candidate, err := QueryCandidate(cliCtx, queryRoute, args[0])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(candidate)
		},
	}
}

func QueryCandidate(cliCtx context.CLIContext, queryRoute, ethAddress string) (candidate types.Candidate, err error) {
	data, err := cliCtx.Codec.MarshalJSON(types.NewQueryCandidateParams(ethAddress))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryCandidate)
	res, err := common.RobustQueryWithData(cliCtx, route, data)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &candidate)
	return
}

// QueryValidators is an interface for convenience to query (all) validators in staking module
func QueryValidators(cliCtx context.CLIContext, storeName string) (validators stakingTypes.Validators, err error) {
	resKVs, _, err := cliCtx.QuerySubspace(stakingTypes.ValidatorsKey, storeName)
	if err != nil {
		return
	}

	for _, kv := range resKVs {
		validators = append(validators, stakingTypes.MustUnmarshalValidator(cliCtx.Codec, kv.Value))
	}
	return
}

// QueryBondedValidators is an interface for convenience to query bonded validators in staking module
func QueryBondedValidators(cliCtx context.CLIContext, storeName string) (validators stakingTypes.Validators, err error) {
	allValidators, err := QueryValidators(cliCtx, storeName)
	if err != nil {
		return
	}

	for _, val := range allValidators {
		if val.Status == sdk.Bonded {
			validators = append(validators, val)
		}
	}

	return
}

// addrStr should be bech32 cosmos account address with prefix cosmos
func QueryValidator(cliCtx context.CLIContext, storeName string, addrStr string) (validator stakingTypes.Validator, err error) {
	addr, err := sdk.AccAddressFromBech32(addrStr)
	if err != nil {
		log.Error(err)
		return
	}

	res, _, err := cliCtx.QueryStore(stakingTypes.GetValidatorKey(sdk.ValAddress(addr)), storeName)
	if err != nil {
		return
	}

	if len(res) == 0 {
		err = fmt.Errorf("No validator found with address %s", addr)
		return
	}

	validator = stakingTypes.MustUnmarshalValidator(cliCtx.Codec, res)
	return
}

// GetCmdReward queries reward info
func GetCmdReward(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reward [ethAddress]",
		Short: "query reward info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			reward, err := QueryReward(cliCtx, queryRoute, args[0])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(reward)
		},
	}
}

// GetCmdRewardRequest queries reward request
func GetCmdRewardRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reward-request [ethAddress]",
		Short: "query reward request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			reward, err := QueryReward(cliCtx, queryRoute, args[0])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			log.Info(string(reward.GetRewardRequest()))
			return nil
		},
	}
}

// Query reward info
func QueryReward(cliCtx context.CLIContext, queryRoute string, ethAddress string) (reward types.Reward, err error) {
	data, err := cliCtx.Codec.MarshalJSON(types.NewQueryRewardParams(ethAddress))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryReward)
	res, err := common.RobustQueryWithData(cliCtx, route, data)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &reward)
	return
}
