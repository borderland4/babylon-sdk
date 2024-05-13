package types

// DefaultParams returns default babylon parameters
func DefaultParams(denom string) Params {
	return Params{ // todo: revisit and set proper defaults
		MaxGasEndBlocker: 500_000,
	}
}

// ValidateBasic performs basic validation on babylon parameters.
func (p Params) ValidateBasic() error {
	if p.MaxGasEndBlocker == 0 {
		return ErrInvalid.Wrap("empty max gas end-blocker setting")
	}
	return nil
}
