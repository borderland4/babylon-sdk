package e2e

import (
	"math/rand"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/babylonchain/babylon/testutil/datagen"
	bbn "github.com/babylonchain/babylon/types"
	bstypes "github.com/babylonchain/babylon/x/btcstaking/types"
	zctypes "github.com/babylonchain/babylon/x/zoneconcierge/types"
	"github.com/btcsuite/btcd/chaincfg"
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
	x := NewExample(t)

	// deploy Babylon contracts to the consumer chain
	consumerCli, consumerContracts, providerCli := x.SetupBabylonIntegration()
	require.NotEmpty(t, consumerCli.Chain.ChainID)
	require.NotEmpty(t, providerCli.Chain.ChainID)
	require.NotEmpty(t, consumerContracts.Babylon)
	require.NotEmpty(t, consumerContracts.BTCStaking)

	// query admin
	adminResp, err := consumerCli.Query(consumerContracts.BTCStaking, Query{"admin": {}})
	require.NoError(t, err)
	require.Equal(t, adminResp["admin"], consumerCli.GetSender().String())

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

	_, err = consumerCli.Exec(consumerContracts.BTCStaking, packetDataBytes)
	require.NoError(t, err)
}
