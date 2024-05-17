package contract

// SudoMsg is a message sent from the Babylon module to a smart contract
// TODO: implement
type SudoMsg struct {
	TestSudoMsg *struct{} `json:"test_sudo_msg,omitempty"`
}
