package e2e

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/babylonchain/babylon/testutil/datagen"
	bbn "github.com/babylonchain/babylon/types"
	bstypes "github.com/babylonchain/babylon/x/btcstaking/types"
	zctypes "github.com/babylonchain/babylon/x/zoneconcierge/types"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

func NewBTCStakingPacketData(packet *bstypes.BTCStakingIBCPacket) *zctypes.ZoneconciergePacketData {
	return &zctypes.ZoneconciergePacketData{
		Packet: &zctypes.ZoneconciergePacketData_BtcStaking{
			BtcStaking: packet,
		},
	}
}

func GenIBCPacket(t *testing.T, r *rand.Rand) *zctypes.ZoneconciergePacketData {

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
	return NewBTCStakingPacketData(packet)
}
