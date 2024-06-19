package contract

import "encoding/hex"

// SudoMsg is a message sent from the Babylon module to a smart contract
type SudoMsg struct {
	BeginBlockMsg struct {
		HashHex    string `json:"hash_hex"`     // HashHex is the hash of the block in hex
		AppHashHex string `json:"app_hash_hex"` // AppHashHex is the app hash of the block in hex
	} `json:"begin_block,omitempty"`
	EndBlockMsg struct {
		HashHex    string `json:"hash_hex"`     // HashHex is the hash of the block in hex
		AppHashHex string `json:"app_hash_hex"` // AppHashHex is the app hash of the block in hex
	} `json:"end_block,omitempty"`
}

func NewBeginBlockSudoMsg(hash []byte, appHash []byte) *SudoMsg {
	msg := struct {
		HashHex    string `json:"hash_hex"`
		AppHashHex string `json:"app_hash_hex"`
	}{hex.EncodeToString(hash), hex.EncodeToString(appHash)}
	return &SudoMsg{
		BeginBlockMsg: msg,
	}
}

func NewEndBlockSudoMsg(hash []byte, appHash []byte) *SudoMsg {
	msg := struct {
		HashHex    string `json:"hash_hex"`
		AppHashHex string `json:"app_hash_hex"`
	}{hex.EncodeToString(hash), hex.EncodeToString(appHash)}
	return &SudoMsg{
		EndBlockMsg: msg,
	}
}
