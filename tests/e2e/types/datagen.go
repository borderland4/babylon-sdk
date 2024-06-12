package types

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/babylonchain/babylon/testutil/datagen"
	bbn "github.com/babylonchain/babylon/types"
	"github.com/babylonchain/babylon/x/btcstaking/types"
	"github.com/btcsuite/btcd/chaincfg"
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

func GenIBCPacket(t *testing.T, r *rand.Rand) ExecuteMessage {
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

	//newFp := &bstypes.NewFinalityProvider{
	//	Description: &ctypes.Description{
	//		Moniker:         "fp1",
	//		Identity:        "Finality Provider 1",
	//		Website:         "https://fp1.com",
	//		SecurityContact: "security_contact",
	//		Details:         "details",
	//	},
	//	Commission: "0.05", // Assuming Decimal::percent(5) converts to "0.05"
	//	BabylonPk:  nil,    // None equivalent in Go is nil
	//	BtcPkHex:   "f1",
	//	Pop: &types.ProofOfPossession{
	//		BtcSigType: 0,
	//		BabylonSig: []byte("mock_pub_rand"),
	//		BtcSig:     []byte("mock_pub_rand"),
	//	},
	//	ConsumerId: "osmosis-1",
	//}

	newFp := NewFinalityProvider{
		Description: &FinalityProviderDescription{
			Moniker:         "fp1",
			Identity:        "Finality Provider 1",
			Website:         "https://fp1.com",
			SecurityContact: "security_contact",
			Details:         "details",
		},
		Commission: "0.05", // Assuming Decimal::percent(5) converts to "0.05"
		BabylonPK: &PubKey{
			Key: base64.StdEncoding.EncodeToString([]byte("mock_pub_rand")),
		}, // None equivalent in Go is nil
		BTCPKHex: "f1",
		Pop: &ProofOfPossession{
			BTCSigType: 0,
			BabylonSig: base64.StdEncoding.EncodeToString([]byte("mock_pub_rand")),
			BTCSig:     base64.StdEncoding.EncodeToString([]byte("mock_pub_rand")),
		},
		ConsumerID: "osmosis-1",
	}

	activeDel := []ActiveBtcDelegation{
		// Add ActiveBtcDelegation instances as needed
	}

	slashedDel := []SlashedBtcDelegation{
		// Add SlashedBtcDelegation instances as needed
	}

	unbondedDel := []UnbondedBtcDelegation{
		// Add UnbondedBtcDelegation instances as needed
	}

	// Create the ExecuteMessage instance
	executeMessage := ExecuteMessage{
		BtcStaking: BtcStaking{
			NewFP:       []NewFinalityProvider{newFp},
			ActiveDel:   activeDel,
			SlashedDel:  slashedDel,
			UnbondedDel: unbondedDel,
		},
	}

	return executeMessage

	//_, mockDel := GenBTCDelegation()
	//activDel, err := CreateActiveBTCDelegation(mockDel)
	//require.NoError(t, err)

	//packet := &bstypes.BTCStakingIBCPacket{
	//	NewFp: []*bstypes.NewFinalityProvider{
	//		newFp,
	//	},
	//	ActiveDel:   []*bstypes.ActiveBTCDelegation{},
	//	SlashedDel:  []*bstypes.SlashedBTCDelegation{},
	//	UnbondedDel: []*bstypes.UnbondedBTCDelegation{},
	//}
	//return NewBTCStakingPacketData(packet)
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

func CreateActiveBTCDelegation(activeDel *bstypes.BTCDelegation) (*bstypes.ActiveBTCDelegation, error) {
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

type NewFinalityProvider struct {
	// Description defines the description terms for the finality provider
	Description *FinalityProviderDescription `json:"description,omitempty"`
	// Commission defines the commission rate of the finality provider
	Commission string `json:"commission"`
	// BabylonPK is the Babylon secp256k1 PK of this finality provider
	BabylonPK *PubKey `json:"babylon_pk,omitempty"`
	// BTCPKHex is the Bitcoin secp256k1 PK of this finality provider
	// the PK follows encoding in BIP-340 spec in hex format
	BTCPKHex string `json:"btc_pk_hex"`
	// PoP is the proof of possession of the babylon_pk and btc_pk
	Pop *ProofOfPossession `json:"pop,omitempty"`
	// ConsumerID is the ID of the consumer that the finality provider is operating on.
	ConsumerID string `json:"consumer_id"`
}

type FinalityProviderDescription struct {
	// Moniker is the name of the finality provider
	Moniker string `json:"moniker"`
	// Identity is the identity of the finality provider
	Identity string `json:"identity"`
	// Website is the website of the finality provider
	Website string `json:"website"`
	// SecurityContact is the security contact of the finality provider
	SecurityContact string `json:"security_contact"`
	// Details is the details of the finality provider
	Details string `json:"details"`
}

type PubKey struct {
	// Key is the compressed public key of the finality provider
	Key string `json:"key"`
}

type ProofOfPossession struct {
	// BTCSigType indicates the type of btc_sig in the pop
	BTCSigType int32 `json:"btc_sig_type"`
	// BabylonSig is the signature generated via sign(sk_babylon, pk_btc)
	BabylonSig string `json:"babylon_sig"`
	// BTCSig is the signature generated via sign(sk_btc, babylon_sig)
	// the signature follows encoding in either BIP-340 spec or BIP-322 spec
	BTCSig string `json:"btc_sig"`
}

// Define the other necessary structs
type ActiveBtcDelegation struct {
	// Define fields as needed
}

type SlashedBtcDelegation struct {
	// Define fields as needed
}

type UnbondedBtcDelegation struct {
	// Define fields as needed
}

// Define the ExecuteMessage struct
type ExecuteMessage struct {
	BtcStaking BtcStaking `json:"btc_staking"`
}

type BtcStaking struct {
	NewFP       []NewFinalityProvider   `json:"new_fp"`
	ActiveDel   []ActiveBtcDelegation   `json:"active_del"`
	SlashedDel  []SlashedBtcDelegation  `json:"slashed_del"`
	UnbondedDel []UnbondedBtcDelegation `json:"unbonded_del"`
}
