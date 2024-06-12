package types

import (
	"encoding/base64"
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
)

func GenExecMessage() ExecuteMessage {
	_, mockDel := GenBTCDelegation()
	ad := ConvertBTCDelegationToActiveBtcDelegation(mockDel)

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
		BTCPKHex: ad.FpBtcPkList[0],
		Pop: &ProofOfPossession{
			BTCSigType: 0,
			BabylonSig: base64.StdEncoding.EncodeToString([]byte("mock_pub_rand")),
			BTCSig:     base64.StdEncoding.EncodeToString([]byte("mock_pub_rand")),
		},
		ConsumerID: "osmosis-1",
	}

	// Create the ExecuteMessage instance
	executeMessage := ExecuteMessage{
		BtcStaking: BtcStaking{
			NewFP:       []NewFinalityProvider{newFp},
			ActiveDel:   []ActiveBtcDelegation{ad},
			SlashedDel:  []SlashedBtcDelegation{},
			UnbondedDel: []UnbondedBtcDelegation{},
		},
	}

	return executeMessage
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
	BTCSigType int32  `json:"btc_sig_type"`
	BabylonSig string `json:"babylon_sig"`
	BTCSig     string `json:"btc_sig"`
}

type CovenantAdaptorSignatures struct {
	CovPK       string   `json:"cov_pk"`       // Public key of the covenant emulator
	AdaptorSigs []string `json:"adaptor_sigs"` // List of adaptor signatures
}

// SignatureInfo represents a signature and its public key
type SignatureInfo struct {
	PK  string `json:"pk"`  // Public key
	Sig string `json:"sig"` // Signature
}

// BtcUndelegationInfo represents the undelegation information
type BtcUndelegationInfo struct {
	UnbondingTx           string                      `json:"unbonding_tx"`                // Unbonding transaction
	DelegatorUnbondingSig string                      `json:"delegator_unbonding_sig"`     // Signature on the unbonding transaction by the delegator
	CovenantUnbondingSigs []SignatureInfo             `json:"covenant_unbonding_sig_list"` // List of signatures on the unbonding transaction by covenant members
	SlashingTx            string                      `json:"slashing_tx"`                 // Unbonding slashing transaction
	DelegatorSlashingSig  string                      `json:"delegator_slashing_sig"`      // Signature on the slashing transaction by the delegator
	CovenantSlashingSigs  []CovenantAdaptorSignatures `json:"covenant_slashing_sigs"`      // List of adaptor signatures on the unbonding slashing transaction by each covenant member
}

type ActiveBtcDelegation struct {
	BTCPkHex             string                      `json:"btc_pk_hex"`             // Bitcoin secp256k1 PK of the BTC delegator in hex format
	FpBtcPkList          []string                    `json:"fp_btc_pk_list"`         // List of BIP-340 PKs of the finality providers
	StartHeight          uint64                      `json:"start_height"`           // Start BTC height of the BTC delegation
	EndHeight            uint64                      `json:"end_height"`             // End BTC height of the BTC delegation
	TotalSat             uint64                      `json:"total_sat"`              // Total BTC stakes in this delegation in satoshi
	StakingTx            string                      `json:"staking_tx"`             // Staking transaction
	SlashingTx           string                      `json:"slashing_tx"`            // Slashing transaction
	DelegatorSlashingSig string                      `json:"delegator_slashing_sig"` // Signature on the slashing transaction by the delegator
	CovenantSigs         []CovenantAdaptorSignatures `json:"covenant_sigs"`          // List of adaptor signatures by covenant members
	StakingOutputIdx     uint32                      `json:"staking_output_idx"`     // Index of the staking output in the staking transaction
	UnbondingTime        uint32                      `json:"unbonding_time"`         // Used in unbonding output time-lock path and slashing transactions change outputs
	UndelegationInfo     *BtcUndelegationInfo        `json:"undelegation_info"`      // Undelegation info of this delegation
	ParamsVersion        uint32                      `json:"params_version"`         // Params version used to validate the delegation
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

func ConvertBTCDelegationToActiveBtcDelegation(mockDel *bstypes.BTCDelegation) ActiveBtcDelegation {
	// Convert the FpBtcPkList from BIP340PubKey to string (assuming a ToHex method exists)
	var fpBtcPkList []string
	for _, pk := range mockDel.FpBtcPkList {
		fpBtcPkList = append(fpBtcPkList, pk.MarshalHex()) // Implement ToHex method for BIP340PubKey
	}

	// Convert CovenantAdaptorSignatures
	var covenantSigs []CovenantAdaptorSignatures
	for _, cs := range mockDel.CovenantSigs {
		var adaptorSigs []string
		for _, sig := range cs.AdaptorSigs {
			adaptorSigs = append(adaptorSigs, base64.StdEncoding.EncodeToString(sig))
		}
		covenantSigs = append(covenantSigs, CovenantAdaptorSignatures{
			CovPK:       cs.CovPk.MarshalHex(),
			AdaptorSigs: adaptorSigs,
		})
	}

	var covenantUnbondingSigs []SignatureInfo
	for _, sigInfo := range mockDel.BtcUndelegation.CovenantUnbondingSigList {
		covenantUnbondingSigs = append(covenantUnbondingSigs, SignatureInfo{
			PK:  sigInfo.Pk.MarshalHex(),
			Sig: base64.StdEncoding.EncodeToString(sigInfo.Sig.MustMarshal()),
		})
	}

	var covenantSlashingSigs []CovenantAdaptorSignatures
	for _, cs := range mockDel.BtcUndelegation.CovenantSlashingSigs {
		var adaptorSigs []string
		for _, sig := range cs.AdaptorSigs {
			adaptorSigs = append(adaptorSigs, base64.StdEncoding.EncodeToString(sig))
		}
		covenantSlashingSigs = append(covenantSlashingSigs, CovenantAdaptorSignatures{
			CovPK:       cs.CovPk.MarshalHex(),
			AdaptorSigs: adaptorSigs,
		})
	}

	// Create and return the ActiveBtcDelegation struct
	return ActiveBtcDelegation{
		BTCPkHex:             mockDel.BtcPk.MarshalHex(), // Implement ToHex method for BIP340PubKey
		FpBtcPkList:          fpBtcPkList,
		StartHeight:          mockDel.StartHeight,
		EndHeight:            mockDel.EndHeight,
		TotalSat:             mockDel.TotalSat,
		StakingTx:            base64.StdEncoding.EncodeToString(mockDel.StakingTx),
		SlashingTx:           base64.StdEncoding.EncodeToString(mockDel.SlashingTx.MustMarshal()),   // Assuming SlashingTx has a TxData field
		DelegatorSlashingSig: base64.StdEncoding.EncodeToString(mockDel.DelegatorSig.MustMarshal()), // Assuming DelegatorSig has a Sig field
		CovenantSigs:         covenantSigs,
		StakingOutputIdx:     mockDel.StakingOutputIdx,
		UnbondingTime:        mockDel.UnbondingTime,
		UndelegationInfo: &BtcUndelegationInfo{
			UnbondingTx:           base64.StdEncoding.EncodeToString(mockDel.BtcUndelegation.UnbondingTx),
			SlashingTx:            base64.StdEncoding.EncodeToString(mockDel.BtcUndelegation.SlashingTx.MustMarshal()),
			DelegatorSlashingSig:  base64.StdEncoding.EncodeToString(mockDel.BtcUndelegation.DelegatorSlashingSig.MustMarshal()),
			CovenantUnbondingSigs: covenantUnbondingSigs,
			CovenantSlashingSigs:  covenantSlashingSigs,
		},
		ParamsVersion: mockDel.ParamsVersion,
	}
}
