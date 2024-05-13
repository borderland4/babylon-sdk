package contract

type (
	CustomMsg struct {
		Test *TestMsg `json:"test,omitempty"`
	}
	TestMsg struct {
		Placeholder string `json:"placeholder,omitempty"`
	}
)
