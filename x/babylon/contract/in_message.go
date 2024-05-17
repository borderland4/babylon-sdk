package contract

// CustomMsg is a message sent from a smart contract to the Babylon module
// TODO: implement
type CustomMsg struct {
	Test *TestMsg `json:"test,omitempty"`
}

type TestMsg struct {
	Placeholder string `json:"placeholder,omitempty"`
}
