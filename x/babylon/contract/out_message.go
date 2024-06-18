package contract

// SudoMsg is a message sent from the Babylon module to a smart contract
type SudoMsg struct {
	BeginBlockMsg *BeginBlock `json:"begin_block,omitempty"`
	EndBlockMsg   *EndBlock   `json:"end_block,omitempty"`
}

type BeginBlock struct {
	HashHex    string `json:"hash_hex"`     // HashHex is the hash of the block in hex
	AppHashHex string `json:"app_hash_hex"` // AppHashHex is the app hash of the block in hex
}

type EndBlock struct {
	HashHex    string `json:"hash_hex"`     // HashHex is the hash of the block in hex
	AppHashHex string `json:"app_hash_hex"` // AppHashHex is the app hash of the block in hex
}
