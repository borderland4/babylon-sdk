package contract

import "time"

// SudoMsg is a message sent from the Babylon module to a smart contract
type SudoMsg struct {
	BeginBlockMsg *BeginBlock `json:"begin_block,omitempty"`
}

type BeginBlock struct {
	Height     int64     `json:"height"`       // Height is the height of the block
	HashHex    string    `json:"hash_hex"`     // HashHex is the hash of the block in hex
	Time       time.Time `json:"time"`         // Time is the time of the block
	ChainID    string    `json:"chain_id"`     // ChainId is the chain ID of the block
	AppHashHex string    `json:"app_hash_hex"` // AppHashHex is the app hash of the block in hex
}
