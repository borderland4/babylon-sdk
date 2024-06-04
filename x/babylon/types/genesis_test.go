package types_test

import (
	"testing"

	"github.com/babylonchain/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateGenesis(t *testing.T) {
	specs := map[string]struct {
		state  types.GenesisState
		expErr bool
	}{
		"default params": {
			state:  *types.DefaultGenesisState(sdk.DefaultBondDenom),
			expErr: false,
		},
		"custom param, should pass": {
			state: types.GenesisState{
				Params: types.Params{
					BabylonContractAddress:    types.EmptyAddr,
					BtcStakingContractAddress: types.EmptyAddr,
					MaxGasBeginBlocker:        600_000,
				},
			},
			expErr: false,
		},
		"custom small value param, should pass": {
			state: types.GenesisState{
				Params: types.Params{
					BabylonContractAddress:    types.EmptyAddr,
					BtcStakingContractAddress: types.EmptyAddr,
					MaxGasBeginBlocker:        10_000,
				},
			},
			expErr: false,
		},
		"invalid max gas length, should fail": {
			state: types.GenesisState{
				Params: types.Params{
					BabylonContractAddress:    types.EmptyAddr,
					BtcStakingContractAddress: types.EmptyAddr,
					MaxGasBeginBlocker:        0,
				},
			},
			expErr: true,
		},
		"invalid max cap coin denom, should fail": {
			state: types.GenesisState{
				Params: types.Params{
					BabylonContractAddress:    types.EmptyAddr,
					BtcStakingContractAddress: types.EmptyAddr,
					MaxGasBeginBlocker:        0,
				},
			},
			expErr: true,
		},
		"invalid max cap coin amount, should fail": {
			state: types.GenesisState{
				Params: types.Params{
					BabylonContractAddress:    types.EmptyAddr,
					BtcStakingContractAddress: types.EmptyAddr,
					MaxGasBeginBlocker:        0,
				},
			},
			expErr: true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			err := types.ValidateGenesis(&spec.state)
			if spec.expErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
