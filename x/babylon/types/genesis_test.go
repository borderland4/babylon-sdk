package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateGenesis(t *testing.T) {
	specs := map[string]struct {
		state  GenesisState
		expErr bool
	}{
		"default params": {
			state:  *DefaultGenesisState(sdk.DefaultBondDenom),
			expErr: false,
		},
		"custom param, should pass": {
			state: GenesisState{
				Params: Params{
					MaxGasEndBlocker: 600_000,
				},
			},
			expErr: false,
		},
		"custom small value param, should pass": {
			state: GenesisState{
				Params: Params{
					MaxGasEndBlocker: 10_000,
				},
			},
			expErr: false,
		},
		"invalid max gas length, should fail": {
			state: GenesisState{
				Params: Params{
					MaxGasEndBlocker: 0,
				},
			},
			expErr: true,
		},
		"invalid max cap coin denom, should fail": {
			state: GenesisState{
				Params: Params{
					MaxGasEndBlocker: 0,
				},
			},
			expErr: true,
		},
		"invalid max cap coin amount, should fail": {
			state: GenesisState{
				Params: Params{
					MaxGasEndBlocker: 0,
				},
			},
			expErr: true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			err := ValidateGenesis(&spec.state)
			if spec.expErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
