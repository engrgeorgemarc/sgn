package multinode

import (
	"math/big"
	"testing"

	"github.com/celer-network/goutils/log"
	tc "github.com/celer-network/sgn/testing/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupValidator(maxValidatorNum *big.Int) {
	log.Infoln("set up new sgn env")
	p := &tc.SGNParams{
		CelrAddr:               tc.E2eProfile.CelrAddr,
		GovernProposalDeposit:  big.NewInt(1), // TODO: use a more practical value
		GovernVoteTimeout:      big.NewInt(1), // TODO: use a more practical value
		BlameTimeout:           big.NewInt(0),
		MinValidatorNum:        big.NewInt(1),
		MaxValidatorNum:        maxValidatorNum,
		MinStakingPool:         big.NewInt(1),
		IncreaseRateWaitTime:   big.NewInt(1), // TODO: use a more practical value
		SidechainGoLiveTimeout: big.NewInt(0),
	}
	tc.SetupNewSGNEnv(p)
	tc.SleepWithLog(10, "sgn being ready")
}

func TestE2EValidator(t *testing.T) {
	t.Run("e2e-validator", func(t *testing.T) {
		t.Run("validatorTest", validatorTest)
		t.Run("replaceValidatorTest", replaceValidatorTest)
	})
}

func validatorTest(t *testing.T) {
	log.Info("===================================================================")
	log.Info("======================== Test validator ===========================")
	setupValidator(big.NewInt(11))

	transactor := tc.NewTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
	)

	// delegation ratio. V0 : V1 : V2 = 2 : 1 : 1
	amts := []*big.Int{big.NewInt(2000000000000000000), big.NewInt(1000000000000000000), big.NewInt(1000000000000000000)}

	// add two validators, 0 and 1
	log.Infoln("---------- It should add two validators successfully ----------")
	for i := 0; i < 2; i++ {
		log.Infoln("Adding validator", i)
		// get auth
		ethAddr, auth, err := tc.GetAuth(tc.ValEthKs[i])
		tc.ChkTestErr(t, err, "failed to get auth")
		tc.AddCandidateWithStake(t, transactor, ethAddr, auth, tc.SgnOperators[i], amts[i], big.NewInt(1), big.NewInt(1), big.NewInt(10000), true)
		tc.CheckValidatorNum(t, transactor, i+1)
	}

	log.Infoln("---------- It should fail to add validator 2 without enough delegation ----------")
	ethAddr, auth, err := tc.GetAuth(tc.ValEthKs[2])
	tc.ChkTestErr(t, err, "failed to get auth")
	initialDelegation := big.NewInt(1)
	tc.AddCandidateWithStake(t, transactor, ethAddr, auth, tc.SgnOperators[2], initialDelegation, big.NewInt(10), big.NewInt(1), big.NewInt(10000), false)
	log.Info("Query sgn about validators to check if validator 2 is not added...")
	tc.CheckValidatorNum(t, transactor, 2)

	log.Infoln("---------- It should correctly add validator 2 with enough delegation ----------")
	err = tc.DelegateStake(auth, ethAddr, big.NewInt(0).Sub(amts[2], initialDelegation))
	tc.ChkTestErr(t, err, "failed to delegate stake")
	tc.CheckValidatorNum(t, transactor, 3)
	tc.CheckValidator(t, transactor, tc.SgnOperators[2], amts[2], sdk.Bonded)

	log.Infoln("---------- It should successfully remove validator 2 caused by intendWithdraw ----------")
	err = tc.IntendWithdraw(auth, ethAddr, amts[2])
	tc.ChkTestErr(t, err, "failed to intendWithdraw stake")
	log.Info("Query sgn about the validators to check if it has correct number of validators...")
	tc.CheckValidatorNum(t, transactor, 2)
	tc.CheckValidatorStatus(t, transactor, tc.SgnOperators[2], sdk.Unbonding)

	err = tc.ConfirmUnbondedCandidate(auth, ethAddr)
	tc.ChkTestErr(t, err, "failed to confirmUnbondedCandidate")
	tc.CheckCandidate(t, transactor, ethAddr, tc.SgnOperators[2], big.NewInt(0))

	err = tc.DelegateStake(auth, ethAddr, amts[2])
	tc.ChkTestErr(t, err, "failed to delegate stake")
	tc.CheckValidatorNum(t, transactor, 3)
	tc.CheckValidator(t, transactor, tc.SgnOperators[2], amts[2], sdk.Bonded)
	// TODO: normally add back validator 1
}

func replaceValidatorTest(t *testing.T) {
	log.Info("===================================================================")
	log.Info("========================  Test replacing validator ===========================")
	setupValidator(big.NewInt(2))

	transactor := tc.NewTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
	)

	amts := []*big.Int{big.NewInt(5000000000000000000), big.NewInt(1000000000000000000), big.NewInt(2000000000000000000)}
	// add two validators, 0 and 1
	tc.AddValidators(t, transactor, tc.ValEthKs[:2], tc.SgnOperators[:2], amts[:2])

	log.Infoln("---------- It should correctly replace validator 1 with validator 2 ----------")
	ethAddr, auth, err := tc.GetAuth(tc.ValEthKs[2])
	tc.ChkTestErr(t, err, "failed to get auth")
	tc.AddCandidateWithStake(t, transactor, ethAddr, auth, tc.SgnOperators[2], amts[2], big.NewInt(1), big.NewInt(1), big.NewInt(10000), true)

	log.Info("Query sgn about the validators...")
	tc.CheckValidatorNum(t, transactor, 2)
	tc.CheckValidator(t, transactor, tc.SgnOperators[1], amts[1], sdk.Unbonding)
}
