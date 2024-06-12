package types

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/babylonchain/babylon/testutil/datagen"
	bbn "github.com/babylonchain/babylon/types"
	"github.com/babylonchain/babylon/x/btcstaking/types"
	"github.com/btcsuite/btcd/chaincfg"
	ctypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	bstypes "github.com/babylonchain/babylon/x/btcstaking/types"
	zctypes "github.com/babylonchain/babylon/x/zoneconcierge/types"
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
	//fpBTCSK, _, err := datagen.GenRandomBTCKeyPair(r)
	//require.NoError(t, err)
	//fpBabylonSK, _, err := datagen.GenRandomSecp256k1KeyPair(r)
	//require.NoError(t, err)
	//fp, err := datagen.GenRandomCustomFinalityProvider(r, fpBTCSK, fpBabylonSK, "consumer-id")
	//require.NoError(t, err)

	//activeDel := &bstypes.ActiveBTCDelegation{
	//	BtcPkHex: fp.BtcPk.MarshalHex(),
	//}
	//fmt.Print(activeDel)

	newFp := &bstypes.NewFinalityProvider{
		Description: &ctypes.Description{
			Moniker:         "fp1",
			Identity:        "Finality Provider 1",
			Website:         "https://fp1.com",
			SecurityContact: "security_contact",
			Details:         "details",
		},
		Commission: "0.05", // Assuming Decimal::percent(5) converts to "0.05"
		BabylonPk:  nil,    // None equivalent in Go is nil
		BtcPkHex:   "f1",
		//Pop: &types.ProofOfPossession{
		//	BtcSigType: 0,
		//	BabylonSig: []byte{},
		//	BtcSig:     []byte{},
		//},
		ConsumerId: "osmosis-1",
	}

	_, mockDel := GenBTCDelegation()
	activDel, err := CreateActiveBTCDelegationEvent(mockDel)
	require.NoError(t, err)

	packet := &bstypes.BTCStakingIBCPacket{
		NewFp: []*bstypes.NewFinalityProvider{
			newFp,
		},
		ActiveDel: []*bstypes.ActiveBTCDelegation{
			activDel,
		},
		SlashedDel:  []*bstypes.SlashedBTCDelegation{},
		UnbondedDel: []*bstypes.UnbondedBTCDelegation{},
	}
	return NewBTCStakingPacketData(packet)
}

var net = &chaincfg.RegressionNetParams

func GenBTCDelegation() (*types.Params, *bstypes.BTCDelegation) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	t := &testing.T{}

	delSK, _, err := datagen.GenRandomBTCKeyPair(r)
	require.NoError(t, err)

	// restaked to a random number of finality providers
	numRestakedFPs := int(datagen.RandomInt(r, 10) + 1)
	_, fpPKs, err := datagen.GenRandomBTCKeyPairs(r, numRestakedFPs)
	require.NoError(t, err)
	fpBTCPKs := bbn.NewBIP340PKsFromBTCPKs(fpPKs)

	// (3, 5) covenant committee
	covenantSKs, covenantPKs, err := datagen.GenRandomBTCKeyPairs(r, 5)
	require.NoError(t, err)
	covenantQuorum := uint32(3)

	stakingTimeBlocks := uint16(5)
	stakingValue := int64(2 * 10e8)
	slashingAddress, err := datagen.GenRandomBTCAddress(r, net)
	require.NoError(t, err)

	slashingRate := sdkmath.LegacyNewDecWithPrec(int64(datagen.RandomInt(r, 41)+10), 2)
	unbondingTime := uint16(100) + 1
	slashingChangeLockTime := unbondingTime

	bsParams := &types.Params{
		CovenantPks:     bbn.NewBIP340PKsFromBTCPKs(covenantPKs),
		CovenantQuorum:  covenantQuorum,
		SlashingAddress: slashingAddress.EncodeAddress(),
	}

	// only the quorum of signers provided the signatures
	covenantSigners := covenantSKs[:covenantQuorum]

	// construct the BTC delegation with everything
	btcDel, err := datagen.GenRandomBTCDelegation(
		r,
		t,
		net,
		fpBTCPKs,
		delSK,
		covenantSigners,
		covenantPKs,
		covenantQuorum,
		slashingAddress.EncodeAddress(),
		1000,
		uint64(1000+stakingTimeBlocks),
		uint64(stakingValue),
		slashingRate,
		slashingChangeLockTime,
	)
	require.NoError(t, err)
	return bsParams, btcDel

	//btcDelBytes, err := btcDel.Marshal()
	//require.NoError(t, err)
	//btcDelPath := filepath.Join(dir, BTC_DEL_FILENAME)
	//err = os.WriteFile(btcDelPath, btcDelBytes, 0644)
	//require.NoError(t, err)

	//paramsBytes, err := bsParams.Marshal()
	//require.NoError(t, err)
	//paramsPath := filepath.Join(dir, BTCSTAKING_PARAMS_FILENAME)
	//err = os.WriteFile(paramsPath, paramsBytes, 0644)
	//require.NoError(t, err)
}

func CreateActiveBTCDelegationEvent(activeDel *bstypes.BTCDelegation) (*bstypes.ActiveBTCDelegation, error) {
	fpBtcPkHexList := make([]string, len(activeDel.FpBtcPkList))
	for i, fpBtcPk := range activeDel.FpBtcPkList {
		fpBtcPkHexList[i] = fpBtcPk.MarshalHex()
	}

	slashingTxBytes, err := activeDel.SlashingTx.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SlashingTx: %w", err)
	}

	delegatorSlashingSigBytes, err := activeDel.DelegatorSig.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DelegatorSig: %w", err)
	}

	if activeDel.BtcUndelegation.DelegatorUnbondingSig != nil {
		return nil, fmt.Errorf("unexpected DelegatorUnbondingSig in active delegation")
	}

	unbondingSlashingTxBytes, err := activeDel.BtcUndelegation.SlashingTx.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal BtcUndelegation.SlashingTx: %w", err)
	}

	delegatorUnbondingSlashingSigBytes, err := activeDel.BtcUndelegation.DelegatorSlashingSig.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal BtcUndelegation.DelegatorSlashingSig: %w", err)
	}

	event := &bstypes.ActiveBTCDelegation{
		BtcPkHex:             activeDel.BtcPk.MarshalHex(),
		FpBtcPkList:          fpBtcPkHexList,
		StartHeight:          activeDel.StartHeight,
		EndHeight:            activeDel.EndHeight,
		TotalSat:             activeDel.TotalSat,
		StakingTx:            activeDel.StakingTx,
		SlashingTx:           slashingTxBytes,
		DelegatorSlashingSig: delegatorSlashingSigBytes,
		CovenantSigs:         activeDel.CovenantSigs,
		StakingOutputIdx:     activeDel.StakingOutputIdx,
		UnbondingTime:        activeDel.UnbondingTime,
		UndelegationInfo: &bstypes.BTCUndelegationInfo{
			UnbondingTx:              activeDel.BtcUndelegation.UnbondingTx,
			SlashingTx:               unbondingSlashingTxBytes,
			DelegatorSlashingSig:     delegatorUnbondingSlashingSigBytes,
			CovenantUnbondingSigList: activeDel.BtcUndelegation.CovenantUnbondingSigList,
			CovenantSlashingSigs:     activeDel.BtcUndelegation.CovenantSlashingSigs,
		},
		ParamsVersion: activeDel.ParamsVersion,
	}

	return event, nil
}
