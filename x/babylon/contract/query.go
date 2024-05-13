package contract

// TODO: revise
type (
	CustomQuery struct {
		Test *TestQuery `json:"test,omitempty"`
	}
	TestQuery struct {
		Placeholder string `json:"placeholder,omitempty"`
	}
	TestResponse struct {
		// MaxCap is the max cap limit
		Placeholder2 string `json:"placeholder2"`
	}
)
