package e2e

import (
	"math/rand"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"github.com/CosmWasm/wasmd/x/wasm/ibctesting"
	"github.com/babylonchain/babylon-sdk/demo/app"
	appparams "github.com/babylonchain/babylon-sdk/demo/app/params"
	"github.com/babylonchain/babylon/testutil/datagen"
	bbn "github.com/babylonchain/babylon/types"
	bstypes "github.com/babylonchain/babylon/x/btcstaking/types"
	zctypes "github.com/babylonchain/babylon/x/zoneconcierge/types"
	"github.com/btcsuite/btcd/chaincfg"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting2 "github.com/cosmos/ibc-go/v8/testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	r = rand.New(rand.NewSource(time.Now().Unix()))
)

// In the Test function, we create and run the suite
func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(BabylonSDKTestSuite))
}

// Define the test suite and include the suite.Suite struct
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
func (suite *BabylonSDKTestSuite) SetupSuite() {
	// overwrite init messages in Babylon
	appparams.SetAddressPrefixes()

	// set up coordinator and chains
	t := suite.T()
	coord := NewIBCCoordinator(t)
	provChain := coord.GetChain(ibctesting2.GetChainID(1))
	consChain := coord.GetChain(ibctesting2.GetChainID(2))

	suite.Coordinator = coord
	suite.ConsumerChain = consChain
	suite.ProviderChain = provChain
	suite.ConsumerApp = consChain.App.(*app.ConsumerApp)
	suite.IbcPath = ibctesting.NewPath(consChain, provChain)
	suite.ProviderDenom = sdk.DefaultBondDenom
	suite.ConsumerDenom = sdk.DefaultBondDenom
	suite.MyProvChainActor = provChain.SenderAccount.GetAddress().String()
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

func (suite *BabylonSDKTestSuite) Test1ContractDeployment() {
	// deploy Babylon contracts to the consumer chain
	consumerCli, consumerContracts, providerCli := suite.SetupBabylonIntegration()
	require.NotEmpty(suite.T(), consumerCli.Chain.ChainID)
	require.NotEmpty(suite.T(), providerCli.Chain.ChainID)
	require.NotEmpty(suite.T(), consumerContracts.Babylon)
	require.NotEmpty(suite.T(), consumerContracts.BTCStaking)

	suite.ProviderCli = providerCli
	suite.ConsumerCli = consumerCli
	suite.ConsumerContract = consumerContracts
}

// TestExample is an example test case
func (suite *BabylonSDKTestSuite) Test2MockFPAndDelegation() {
	t := suite.T()

	// query admin
	adminResp, err := suite.ConsumerCli.Query(suite.ConsumerContract.BTCStaking, Query{"admin": {}})
	require.NoError(t, err)
	require.Equal(t, adminResp["admin"], suite.ConsumerCli.GetSender().String())

	// generate a finality provider
	fpBTCSK, _, err := datagen.GenRandomBTCKeyPair(r)
	require.NoError(t, err)
	fpBabylonSK, _, err := datagen.GenRandomSecp256k1KeyPair(r)
	require.NoError(t, err)
	fp, err := datagen.GenRandomCustomFinalityProvider(r, fpBTCSK, fpBabylonSK, "consumer-id")
	require.NoError(t, err)

	// generate a BTC delegation
	delSK, _, err := datagen.GenRandomBTCKeyPair(r)
	require.NoError(t, err)
	covenantSKs, covenantPKs, covenantQuorum := datagen.GenCovenantCommittee(r)
	slashingAddress, err := datagen.GenRandomBTCAddress(r, &chaincfg.RegressionNetParams)
	require.NoError(t, err)
	slashingRate := sdkmath.LegacyNewDecWithPrec(int64(datagen.RandomInt(r, 41)+10), 2)
	slashingChangeLockTime := uint16(101)
	del, err := datagen.GenRandomBTCDelegation(
		r,
		t,
		&chaincfg.RegressionNetParams,
		[]bbn.BIP340PubKey{*fp.BtcPk},
		delSK,
		covenantSKs,
		covenantPKs,
		covenantQuorum,
		slashingAddress.EncodeAddress(),
		1, 1000, 10000,
		slashingRate,
		slashingChangeLockTime,
	)
	require.NoError(t, err)

	packet := &bstypes.BTCStakingIBCPacket{
		NewFp: []*bstypes.NewFinalityProvider{
			// TODO: fill empty data
			&bstypes.NewFinalityProvider{
				// Description: fp.Description,
				Commission: fp.Commission.String(),
				// BabylonPk:  fp.BabylonPk,
				BtcPkHex: fp.BtcPk.MarshalHex(),
				// Pop:        fp.Pop,
				ConsumerId: fp.ConsumerId,
			},
		},
		ActiveDel: []*bstypes.ActiveBTCDelegation{
			&bstypes.ActiveBTCDelegation{
				BtcPkHex:             del.BtcPk.MarshalHex(),
				FpBtcPkList:          []string{del.FpBtcPkList[0].MarshalHex()},
				StartHeight:          del.StartHeight,
				EndHeight:            del.EndHeight,
				TotalSat:             del.TotalSat,
				StakingTx:            del.StakingTx,
				SlashingTx:           *del.SlashingTx,
				DelegatorSlashingSig: *del.DelegatorSig,
				CovenantSigs:         del.CovenantSigs,
				UnbondingTime:        del.UnbondingTime,
				UndelegationInfo: &bstypes.BTCUndelegationInfo{
					UnbondingTx:              del.BtcUndelegation.UnbondingTx,
					CovenantUnbondingSigList: del.BtcUndelegation.CovenantUnbondingSigList,
					SlashingTx:               *del.BtcUndelegation.SlashingTx,
					DelegatorSlashingSig:     *del.BtcUndelegation.DelegatorSlashingSig,
					CovenantSlashingSigs:     del.BtcUndelegation.CovenantSlashingSigs,
				},
				ParamsVersion: del.ParamsVersion,
			},
		},
		SlashedDel:  []*bstypes.SlashedBTCDelegation{},
		UnbondedDel: []*bstypes.UnbondedBTCDelegation{},
	}
	packetData := NewBTCStakingPacketData(packet)

	packetDataBytes, err := zctypes.ModuleCdc.MarshalJSON(packetData)
	require.NoError(t, err)

	_, err = suite.ConsumerCli.Exec(suite.ConsumerContract.BTCStaking, packetDataBytes)
	require.NoError(t, err)
}

// TearDownSuite runs once after all the suite's tests have been run
func (suite *BabylonSDKTestSuite) TearDownSuite() {
	// Cleanup code here
}
