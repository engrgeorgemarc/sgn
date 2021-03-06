package cli

import (
	"github.com/celer-network/sgn/x/slash/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	slashTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "slash transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	slashTxCmd.AddCommand(flags.PostCommands()...)

	return slashTxCmd
}
