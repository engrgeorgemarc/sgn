package testing

import (
	"context"
	"math/big"
	"strings"

	"github.com/celer-network/cChannel-eth-go/deploy"
	"github.com/celer-network/cChannel-eth-go/ledger"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeployLedgerContract() mainchain.Addr {
	ctx := context.Background()
	channelAddrBundle := deploy.DeployAll(DefaultTestEthClient.Auth, DefaultTestEthClient.Client, ctx, 0)
	ledgerAddr := channelAddrBundle.CelerLedgerAddr

	// Disable channel deposit limit
	LogBlkNum(DefaultTestEthClient.Client)
	ledgerContract, err := ledger.NewCelerLedger(ledgerAddr, DefaultTestEthClient.Client)
	ChkErr(err, "failed to NewCelerLedger")
	tx, err := ledgerContract.DisableBalanceLimits(DefaultTestEthClient.Auth)
	ChkErr(err, "failed disable channel deposit limits")
	WaitMinedWithChk(ctx, DefaultTestEthClient.Client, tx, 0, "Disable balance limit")

	log.Infoln("Ledger address:", ledgerAddr.String())
	return ledgerAddr
}

func DeployERC20Contract() (mainchain.Addr, *mainchain.ERC20) {
	ctx := context.Background()
	initAmt := new(big.Int)
	initAmt.SetString("5"+strings.Repeat("0", 44), 10)
	erc20Addr, tx, erc20, err := mainchain.DeployERC20(DefaultTestEthClient.Auth, DefaultTestEthClient.Client, initAmt, "Celer", 18, "CELR")
	ChkErr(err, "failed to deploy ERC20")
	WaitMinedWithChk(ctx, DefaultTestEthClient.Client, tx, 0, "Deploy ERC20 "+erc20Addr.Hex())

	log.Infoln("Erc20 address:", erc20Addr.String())
	return erc20Addr, erc20
}

func DeployGuardContract(sgnParams *SGNParams) mainchain.Addr {
	ctx := context.Background()
	guardAddr, tx, _, err := mainchain.DeployGuard(DefaultTestEthClient.Auth, DefaultTestEthClient.Client, sgnParams.CelrAddr, sgnParams.BlameTimeout, sgnParams.MinValidatorNum, sgnParams.MinStakingPool, sgnParams.SidechainGoLiveTimeout)
	ChkErr(err, "failed to deploy Guard contract")
	WaitMinedWithChk(ctx, DefaultTestEthClient.Client, tx, 0, "Deploy Guard "+guardAddr.Hex())

	log.Infoln("Guard address:", guardAddr.String())
	return guardAddr
}

func DeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy contracts",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			viper.SetConfigFile(viper.GetString(common.FlagConfig))
			err = viper.ReadInConfig()
			if err != nil {
				return
			}

			err = DefaultTestEthClient.SetClient(viper.GetString(common.FlagEthWS))
			if err != nil {
				return
			}

			err = DefaultTestEthClient.SetAuth(viper.GetString(common.FlagEthKeystore), viper.GetString(common.FlagEthPassphrase))
			if err != nil {
				return
			}

			SetEthBaseKs("./docker-volumes/geth-env")
			err = FundAddr("1"+strings.Repeat("0", 20), []*mainchain.Addr{&Client0Addr, &Client1Addr})
			ChkErr(err, "fund client0 and client1")

			ledgerAddr := DeployLedgerContract()
			viper.Set(common.FlagEthLedgerAddress, ledgerAddr.String())

			erc20Addr, erc20 := DeployERC20Contract()
			sgnParams := &SGNParams{
				BlameTimeout:           big.NewInt(50),
				MinValidatorNum:        big.NewInt(0),
				MinStakingPool:         big.NewInt(100),
				SidechainGoLiveTimeout: big.NewInt(0),
				CelrAddr:               erc20Addr,
			}
			guardAddr := DeployGuardContract(sgnParams)
			viper.Set(common.FlagEthGuardAddress, guardAddr.String())
			viper.WriteConfig()

			amt := new(big.Int)
			amt.SetString("1"+strings.Repeat("0", 19), 10)
			tx, err := erc20.Approve(DefaultTestEthClient.Auth, guardAddr, amt)
			ChkErr(err, "failed to deploy Guard contract")
			WaitMinedWithChk(context.Background(), DefaultTestEthClient.Client, tx, 0, "approve erc20")
			return
		},
	}

	return cmd
}