package contract

import "time"

// SudoMsg is a message sent from the Babylon module to a smart contract
type SudoMsg struct {
	BeginBlockMsg *BeginBlock `json:"begin_block,omitempty"`
}

type BeginBlock struct {
	Height     int64     `json:"height"`       // Height returns the height of the block
	HashHex    string    `json:"hash_hex"`     // HashHex returns the hash of the block header in hex
	Time       time.Time `json:"time"`         // Time returns the time of the block
	ChainID    string    `json:"chain_id"`     // ChainId returns the chain ID of the block
	AppHashHex string    `json:"app_hash_hex"` // AppHashHex used in the current block header in hex
}
