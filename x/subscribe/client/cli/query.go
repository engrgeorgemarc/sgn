package cli

import (
	"fmt"

	"github.com/celer-network/sgn/x/subscribe/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	subscribeQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the subscribe module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	subscribeQueryCmd.AddCommand(client.GetCommands(
		GetCmdSubscription(storeKey, cdc),
		GetCmdRequest(storeKey, cdc),
	)...)
	return subscribeQueryCmd
}

// GetCmdSubscription queries subscription info
func GetCmdSubscription(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "subscription [ethAddress]",
		Short: "query subscription info associated with the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			subscription, err := QuerySubscrption(cdc, cliCtx, queryRoute, args[0])
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(subscription)
		},
	}
}

// GetCmdRequest queries request info
func GetCmdRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "request [channelId]",
		Short: "query request info associated with the channelId",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			request, err := QueryRequest(cdc, cliCtx, queryRoute, ethcommon.Hex2Bytes(args[0]))
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(request)
		},
	}
}

func QuerySubscrption(cdc *codec.Codec, cliCtx context.CLIContext, queryRoute, ethAddress string) (subscription types.Subscription, err error) {
	data, err := cdc.MarshalJSON(types.NewQuerySubscrptionParams(ethAddress))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySubscrption)
	res, _, err := cliCtx.QueryWithData(route, data)
	if err != nil {
		fmt.Printf("query error", err)
		return
	}

	cdc.MustUnmarshalJSON(res, &subscription)
	return
}

// Query request info
func QueryRequest(cdc *codec.Codec, cliCtx context.CLIContext, queryRoute string, channelId []byte) (request types.Request, err error) {
	data, err := cdc.MarshalJSON(types.NewQueryRequestParams(channelId))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequest)
	res, _, err := cliCtx.QueryWithData(route, data)
	if err != nil {
		fmt.Printf("query error", err)
		return
	}

	cdc.MustUnmarshalJSON(res, &request)
	return
}
