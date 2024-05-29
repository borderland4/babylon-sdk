package e2e

import (
	"math/rand"
	"testing"
	"time"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/babylonchain/babylon/testutil/datagen"
	bstypes "github.com/babylonchain/babylon/x/btcstaking/types"
	zctypes "github.com/babylonchain/babylon/x/zoneconcierge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

func NewBTCStakingPacketData(packet *bstypes.BTCStakingIBCPacket) *zctypes.ZoneconciergePacketData {
	return &zctypes.ZoneconciergePacketData{
		Packet: &zctypes.ZoneconciergePacketData_BtcStaking{
			BtcStaking: packet,
		},
	}
}

func TestMVP(t *testing.T) {
	// create a provider chain and a consumer chain
	x := setupExampleChains(t)

	// deploy Babylon contracts to the consumer chain
	consumerCli, consumerContracts, providerCli := setupBabylonIntegration(t, x)
	require.NotEmpty(t, consumerCli.Chain.ChainID)
	require.NotEmpty(t, providerCli.Chain.ChainID)
	require.False(t, consumerContracts.Babylon.Empty())
	require.False(t, consumerContracts.BTCStaking.Empty())

	// inject some finality providers via admin commands
	fpBTCSK, _, err := datagen.GenRandomBTCKeyPair(r)
	require.NoError(t, err)
	fpBabylonSK, _, err := datagen.GenRandomSecp256k1KeyPair(r)
	require.NoError(t, err)
	fp, err := datagen.GenRandomCustomFinalityProvider(r, fpBTCSK, fpBabylonSK, "consumer-id")
	require.NoError(t, err)
	fp.Description.Identity = "doesntmatter"
	fp.Description.Website = "website"
	fp.Description.SecurityContact = "website"
	fp.Description.Details = "website"

	packet := &bstypes.BTCStakingIBCPacket{
		NewFp: []*bstypes.NewFinalityProvider{
			&bstypes.NewFinalityProvider{
				Description: fp.Description,
				Commission:  fp.Commission.String(),
				// BabylonPk:   fp.BabylonPk, // TODO: figure out why
				BtcPkHex: fp.BtcPk.MarshalHex(),
				// Pop:        fp.Pop,
				ConsumerId: fp.ConsumerId,
			},
		},
	}
	packetData := NewBTCStakingPacketData(packet)

	newFPPacketBytes, err := zctypes.ModuleCdc.MarshalJSON(packetData)
	require.NoError(t, err)
	msg := &wasmtypes.MsgExecuteContract{
		Sender:   consumerCli.Chain.SenderAccount.GetAddress().String(),
		Contract: consumerContracts.BTCStaking.String(),
		Msg:      newFPPacketBytes,
		Funds:    sdk.NewCoins(),
	}

	// res, err := consumerCli.Chain.SendMsgs(msg)
	// require.NoError(t, err, res)

	wasmKeeper := consumerCli.Chain.App.GetWasmKeeper()
	ms := wasmkeeper.NewMsgServerImpl(&wasmKeeper)
	resp, err := ms.ExecuteContract(consumerCli.Chain.GetContext(), msg)
	require.NoError(t, err, resp)

	// inject some BTC delegations via admin commands

	//
}
