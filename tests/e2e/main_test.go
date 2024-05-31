package e2e

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/CosmWasm/wasmd/x/wasm/ibctesting"
	"github.com/babylonchain/babylon-sdk/demo/app"
	appparams "github.com/babylonchain/babylon-sdk/demo/app/params"
	"github.com/babylonchain/babylon-sdk/tests/e2e/types"
	zctypes "github.com/babylonchain/babylon/x/zoneconcierge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting2 "github.com/cosmos/ibc-go/v8/testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

// In the Test function, we create and run the suite
func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(BabylonSDKTestSuite))
}

// Define the test suite and include the s.Suite struct
type BabylonSDKTestSuite struct {
	suite.Suite

	Coordinator      *ibctesting.Coordinator
	ConsumerChain    *ibctesting.TestChain
	ProviderChain    *ibctesting.TestChain
	ConsumerApp      *app.ConsumerApp
	IbcPath          *ibctesting.Path
	ProviderDenom    string
	ConsumerDenom    string
	MyProvChainActor string

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

func (x *BabylonSDKTestSuite) SetupBabylonIntegration() (*TestConsumerClient, *ConsumerContract, *TestProviderClient) {
	x.Coordinator.SetupConnections(x.IbcPath)

	// setup contracts on consumer
	consumerCli := NewConsumerClient(x.T(), x.ConsumerChain)
	consumerContracts, err := consumerCli.BootstrapContracts()
	require.NoError(x.T(), err)
	// consumerPortID := wasmkeeper.PortIDForContract(consumerContracts.Babylon)

	// add some fees so that we can distribute something
	x.ConsumerChain.DefaultMsgFees = sdk.NewCoins(sdk.NewCoin(x.ConsumerDenom, math.NewInt(1_000_000)))

	providerCli := NewProviderClient(x.T(), x.ProviderChain)

	return consumerCli, consumerContracts, providerCli

	// TODO: fix IBC channel below
	// // setup ibc path
	// x.IbcPath.EndpointB.ChannelConfig = &ibctesting2.ChannelConfig{
	// 	PortID: "zoneconcierge", // TODO: replace this chain/port with Babylon
	// 	Order:  types2.ORDERED,
	// }
	// x.IbcPath.EndpointA.ChannelConfig = &ibctesting2.ChannelConfig{
	// 	PortID: consumerPortID,
	// 	Order:  types2.ORDERED,
	// }
	// x.Coordinator.CreateChannels(x.IbcPath)

	// // when ibc package is relayed
	// require.NotEmpty(x.T(), x.ConsumerChain.PendingSendPackets)
	// require.NoError(x.T(), x.Coordinator.RelayAndAckPendingPackets(x.IbcPath))

	// return consumerCli, consumerContracts, providerCli
}

func (s *BabylonSDKTestSuite) Test1ContractDeployment() {
	// deploy Babylon contracts to the consumer chain
	consumerCli, consumerContracts, providerCli := s.SetupBabylonIntegration()
	require.NotEmpty(s.T(), consumerCli.Chain.ChainID)
	require.NotEmpty(s.T(), providerCli.Chain.ChainID)
	require.NotEmpty(s.T(), consumerContracts.Babylon)
	require.NotEmpty(s.T(), consumerContracts.BTCStaking)

	s.ProviderCli = providerCli
	s.ConsumerCli = consumerCli
	s.ConsumerContract = consumerContracts

	// query admin
	adminResp, err := s.ConsumerCli.Query(s.ConsumerContract.BTCStaking, Query{"admin": {}})
	s.NoError(err)
	s.Equal(adminResp["admin"], s.ConsumerCli.GetSender().String())
}

// TestExample is an example test case
func (s *BabylonSDKTestSuite) Test2MockFPAndDelegation() {
	t := s.T()

	packet := types.GenIBCPacket(t, r)
	packetBytes, err := zctypes.ModuleCdc.MarshalJSON(packet)
	require.NoError(t, err)
	fmt.Println(string(packetBytes))

	_, err = s.ConsumerCli.Exec(s.ConsumerContract.BTCStaking, packetBytes)
	require.NoError(t, err)
}

// TearDownSuite runs once after all the suite's tests have been run
func (s *BabylonSDKTestSuite) TearDownSuite() {
	// Cleanup code here
}
