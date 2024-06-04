package keeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonchain/babylon-sdk/x/babylon/types"
)

func TestInitGenesis(t *testing.T) {
	specs := map[string]struct {
		state  types.GenesisState
		expErr bool
	}{
		"custom param, should pass": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker: 600_000,
				},
			},
			expErr: false,
		},
		"custom small value param, should pass": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker: 10_000,
				},
			},
			expErr: false,
		},
	}

	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			pCtx, keepers := NewTestKeepers(t)
			k := keepers.BabylonKeeper

			k.InitGenesis(pCtx, spec.state)

			p := k.GetParams(pCtx)
			assert.Equal(t, spec.state.Params.MaxGasBeginBlocker, p.MaxGasBeginBlocker)
		})
	}
}

func TestExportGenesis(t *testing.T) {
	pCtx, keepers := NewTestKeepers(t)
	k := keepers.BabylonKeeper
	params := types.DefaultParams(sdk.DefaultBondDenom)

	err := k.SetParams(pCtx, params)
	require.NoError(t, err)

	exported := k.ExportGenesis(pCtx)
	assert.Equal(t, params.MaxGasBeginBlocker, exported.Params.MaxGasBeginBlocker)
}
