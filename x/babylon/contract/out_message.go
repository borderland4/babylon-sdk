package contract

import "time"

// SudoMsg is a message sent from the Babylon module to a smart contract
type SudoMsg struct {
	BeginBlock *BeginBlockMsg `json:"begin_block_msg,omitempty"`
}

type BeginBlockMsg struct {
	Height  int64     `json:"height"`   // Height returns the height of the block
	Hash    []byte    `json:"hash"`     // Hash returns the hash of the block header
	Time    time.Time `json:"time"`     // Time returns the time of the block
	ChainID string    `json:"chain_id"` // ChainId returns the chain ID of the block
	AppHash []byte    `json:"app_hash"` // AppHash used in the current block header
}
