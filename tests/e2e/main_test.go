package e2e

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm/ibctesting"
	"github.com/babylonchain/babylon-sdk/demo/app"
	appparams "github.com/babylonchain/babylon-sdk/demo/app/params"
	"github.com/babylonchain/babylon-sdk/tests/e2e/types"
	bbntypes "github.com/babylonchain/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting2 "github.com/cosmos/ibc-go/v8/testing"
	"github.com/stretchr/testify/suite"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

// In the Test function, we create and run the suite
func TestBabylonSDKTestSuite(t *testing.T) {
	suite.Run(t, new(BabylonSDKTestSuite))
}

// Define the test suite and include the s.Suite struct
type BabylonSDKTestSuite struct {
	suite.Suite

	// provider/consumer and their metadata
	Coordinator      *ibctesting.Coordinator
	ConsumerChain    *ibctesting.TestChain
	ProviderChain    *ibctesting.TestChain
	ConsumerApp      *app.ConsumerApp
	IbcPath          *ibctesting.Path
	ProviderDenom    string
	ConsumerDenom    string
	MyProvChainActor string

	// clients side information
	ProviderCli      *TestProviderClient
	ConsumerCli      *TestConsumerClient
	ConsumerContract *ConsumerContract
}

// SetupSuite runs once before the suite's tests are run
func (s *BabylonSDKTestSuite) SetupSuite() {
	// overwrite init messages in Babylon
	appparams.SetAddressPrefixes()

	// set up coordinator and chains
	t := s.T()
	coord := NewIBCCoordinator(t)
	provChain := coord.GetChain(ibctesting2.GetChainID(1))
	consChain := coord.GetChain(ibctesting2.GetChainID(2))

	s.Coordinator = coord
	s.ConsumerChain = consChain
	s.ProviderChain = provChain
	s.ConsumerApp = consChain.App.(*app.ConsumerApp)
	s.IbcPath = ibctesting.NewPath(consChain, provChain)
	s.ProviderDenom = sdk.DefaultBondDenom
	s.ConsumerDenom = sdk.DefaultBondDenom
	s.MyProvChainActor = provChain.SenderAccount.GetAddress().String()
}

func (x *BabylonSDKTestSuite) setupBabylonIntegration() (*TestConsumerClient, *ConsumerContract, *TestProviderClient) {
	x.Coordinator.SetupConnections(x.IbcPath)

	// consumer client
	consumerCli := NewConsumerClient(x.T(), x.ConsumerChain)
	// setup contracts on consumer
	consumerContracts, err := consumerCli.BootstrapContracts()
	x.NoError(err)
	// provider client
	providerCli := NewProviderClient(x.T(), x.ProviderChain)

	return consumerCli, consumerContracts, providerCli
}

func (s *BabylonSDKTestSuite) Test1ContractDeployment() {
	// deploy Babylon contracts to the consumer chain
	consumerCli, consumerContracts, providerCli := s.setupBabylonIntegration()
	s.NotEmpty(consumerCli.Chain.ChainID)
	s.NotEmpty(providerCli.Chain.ChainID)
	s.NotEmpty(consumerContracts.Babylon)
	s.NotEmpty(consumerContracts.BTCStaking)

	s.ProviderCli = providerCli
	s.ConsumerCli = consumerCli
	s.ConsumerContract = consumerContracts

	// query admin
	adminResp, err := s.ConsumerCli.Query(s.ConsumerContract.BTCStaking, Query{"admin": {}})
	s.NoError(err)
	s.Equal(adminResp["admin"], s.ConsumerCli.GetSender().String())

	// update the contract address in parameters
	msgUpdateParams := &bbntypes.MsgUpdateParams{
		Authority: s.ConsumerApp.BabylonKeeper.GetAuthority(),
		Params: bbntypes.Params{
			MaxGasBeginBlocker:        500_000,
			BabylonContractAddress:    s.ConsumerContract.Babylon.String(),
			BtcStakingContractAddress: s.ConsumerContract.BTCStaking.String(),
		},
	}
	s.ConsumerCli.MustExecGovProposal(msgUpdateParams)
}

// TestExample is an example test case
func (s *BabylonSDKTestSuite) Test2MockConsumerFpDelegation() {
	msg := types.GenExecMessage()
	msgBytes, err := json.Marshal(msg)
	s.NoError(err)

	// send msg to BTC staking contract via admin account
	_, err = s.ConsumerCli.Exec(s.ConsumerContract.BTCStaking, msgBytes)
	s.NoError(err)

	// ensure the finality provider is on consumer chain
	consumerFps, err := s.ConsumerCli.Query(s.ConsumerContract.BTCStaking, Query{"finality_providers": {}})
	s.NoError(err)
	s.NotEmpty(consumerFps)

	// ensure delegations are on consumer chain
	consumerDels, err := s.ConsumerCli.Query(s.ConsumerContract.BTCStaking, Query{"delegations": {}})
	s.NoError(err)
	s.NotEmpty(consumerDels)
}

// TODO: trigger BeginBlock via s.ConsumerChain rather than ConsumerApp
func (s *BabylonSDKTestSuite) Test3BeginBlock() {
	err := s.ConsumerApp.BabylonKeeper.BeginBlocker(s.ConsumerChain.GetContext())
	s.NoError(err)
}

// TearDownSuite runs once after all the suite's tests have been run
func (s *BabylonSDKTestSuite) TearDownSuite() {
	// Cleanup code here
}
