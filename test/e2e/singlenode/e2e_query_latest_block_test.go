package singlenode

import (
	"context"
	"os"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/global"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setUpQueryLatestBlock() []tf.Killable {
	res := setupNewSGNEnv(nil, "query_latest_block")
	tf.SleepWithLog(10, "sgn syncing")

	return res
}

func TestE2EQueryLatestBlock(t *testing.T) {
	toKill := setUpQueryLatestBlock()
	defer tf.TearDown(toKill)

	t.Run("e2e-queryLatestBlock", func(t *testing.T) {
		t.Run("queryLatestBlockTest", queryLatestBlockTest)
	})
}

func queryLatestBlockTest(t *testing.T) {
	// t.Parallel()

	log.Info("=====================================================================")
	log.Info("======================== Test queryLatestBlock ===========================")

	conn, err := ethclient.Dial(tf.EthInstance)
	if err != nil {
		os.Exit(1)
	}

	transactor := tf.NewTransactor(
		t,
		viper.GetString(common.FlagCLIHome),
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetStringSlice(common.FlagSgnTransactors)[0],
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
	)

	blockSGN, err := global.CLIQueryLatestBlock(transactor.CliCtx, global.RouterKey)
	tf.ChkTestErr(t, err, "failed to query latest synced block on sgn")
	log.Infof("Latest block number on SGN is %d", blockSGN.Number)

	header, err := conn.HeaderByNumber(context.Background(), nil)
	tf.ChkTestErr(t, err, "failed to query latest synced block on mainchain")
	log.Infof("Latest block number on mainchain is %d", header.Number)

	assert.GreaterOrEqual(t, header.Number.Uint64(), blockSGN.Number, "blkNumMain should be greater than or equal to blockSGN.Number")
	assert.Greater(t, blockSGN.Number, uint64(0), "blockSGN.Number should be larger than 0")
}