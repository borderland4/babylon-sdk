package types

import (
	"github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// DefaultParams returns default babylon parameters
func DefaultParams(denom string) Params {
	return Params{
		BabylonContractAddress:    sdk.AccAddress(rand.Bytes(address.Len)).String(),
		BtcStakingContractAddress: sdk.AccAddress(rand.Bytes(address.Len)).String(),
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
