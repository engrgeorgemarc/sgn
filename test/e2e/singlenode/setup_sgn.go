package singlenode

import (
	"context"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

func setupNewSGNEnv(sgnParams *tc.SGNParams, testName string) []tc.Killable {
	if sgnParams == nil {
		sgnParams = &tc.SGNParams{
			CelrAddr:               tc.E2eProfile.CelrAddr,
			GovernProposalDeposit:  big.NewInt(1), // TODO: use a more practical value
			GovernVoteTimeout:      big.NewInt(1), // TODO: use a more practical value
			BlameTimeout:           big.NewInt(50),
			MinValidatorNum:        big.NewInt(1),
			MaxValidatorNum:        big.NewInt(11),
			MinStakingPool:         big.NewInt(100),
			IncreaseRateWaitTime:   big.NewInt(1), // TODO: use a more practical value
			SidechainGoLiveTimeout: big.NewInt(0),
		}
	}
	var tx *types.Transaction
	tx, tc.E2eProfile.DPoSAddr, tc.E2eProfile.SGNAddr = tc.DeployDPoSSGNContracts(sgnParams)
	tc.WaitMinedWithChk(context.Background(), tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "DeployDPoSSGNContracts")

	updateSGNConfig()

	sgnProc, err := startSidechain(outRootDir, testName)
	tc.ChkErr(err, "start sidechain")
	tc.SetContracts(tc.E2eProfile.DPoSAddr, tc.E2eProfile.SGNAddr, tc.E2eProfile.LedgerAddr)

	killable := []tc.Killable{sgnProc}
	if sgnParams.StartGateway {
		gatewayProc, err := StartGateway(outRootDir, testName)
		tc.ChkErr(err, "start gateway")
		killable = append(killable, gatewayProc)
	}

	return killable
}

func updateSGNConfig() {
	log.Infoln("Updating SGN's config.json")

	viper.SetConfigFile("../../../config.json")
	err := viper.ReadInConfig()
	tc.ChkErr(err, "failed to read config")

	clientKeystore, err := filepath.Abs("../../keys/ethks0.json")
	tc.ChkErr(err, "get client keystore path")

	viper.Set(common.FlagEthGateway, tc.LocalGeth)
	viper.Set(common.FlagEthDPoSAddress, tc.E2eProfile.DPoSAddr)
	viper.Set(common.FlagEthSGNAddress, tc.E2eProfile.SGNAddr)
	viper.Set(common.FlagEthLedgerAddress, tc.E2eProfile.LedgerAddr)
	viper.Set(common.FlagEthKeystore, clientKeystore)
	viper.WriteConfig()
}

func installSgn() error {
	cmd := exec.Command("make", "install")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "WITH_CLEVELDB=yes")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../../..")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("cp", "./test/data/local_config.json", "./config.json")
	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../../..")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// startSidechain starts sgn sidechain with the data in test/data
func startSidechain(rootDir, testName string) (*os.Process, error) {
	cmd := exec.Command("make", "update-test-data")
	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../../..")
	if err := cmd.Run(); err != nil {
		log.Errorln("Failed to run \"make update-test-data\": ", err)
		return nil, err
	}

	cmd = exec.Command("sgnd", "start")
	cmd.Dir, _ = filepath.Abs("../../..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Errorln("Failed to run \"sgnd start\": ", err)
		return nil, err
	}

	log.Infoln("sgn pid:", cmd.Process.Pid)
	return cmd.Process, nil
}

func StartGateway(rootDir, testName string) (*os.Process, error) {
	cmd := exec.Command("sgncli", "gateway")
	cmd.Dir, _ = filepath.Abs("../../..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	log.Infoln("gateway pid:", cmd.Process.Pid)
	return cmd.Process, nil
}
