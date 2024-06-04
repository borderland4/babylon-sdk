package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var EmptyAddr = sdk.AccAddress([]byte{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}).String()

// DefaultParams returns default babylon parameters
func DefaultParams(denom string) Params {
	return Params{
		BabylonContractAddress:    EmptyAddr,
		BtcStakingContractAddress: EmptyAddr,
		MaxGasBeginBlocker:        500_000,
	}
}

// ValidateBasic performs basic validation on babylon parameters.
func (p Params) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(p.BabylonContractAddress); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(p.BtcStakingContractAddress); err != nil {
		return err
	}
	if p.MaxGasBeginBlocker == 0 {
		return ErrInvalid.Wrap("empty max gas end-blocker setting")
	}
	return nil
}
