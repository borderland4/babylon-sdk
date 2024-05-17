package contract

// CustomQuery is a query request from a smart contract to the Babylon module
// TODO: implement
type CustomQuery struct {
	Test *TestQuery `json:"test,omitempty"`
}

type TestQuery struct {
	Placeholder string `json:"placeholder,omitempty"`
}

type TestResponse struct {
	// MaxCap is the max cap limit
	Placeholder2 string `json:"placeholder2"`
}
