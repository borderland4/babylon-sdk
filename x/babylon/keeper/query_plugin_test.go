package keeper

import (
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	"github.com/cometbft/cometbft/libs/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestChainedCustomQuerier(t *testing.T) {
	myContractAddr := sdk.AccAddress(rand.Bytes(32))
	pCtx, keepers := CreateDefaultTestInput(t)

	specs := map[string]struct {
		src           wasmvmtypes.QueryRequest
		viewKeeper    viewKeeper
		expData       []byte
		expErr        bool
		expNextCalled bool
	}{
		"non custom query": {
			src: wasmvmtypes.QueryRequest{
				Bank: &wasmvmtypes.BankQuery{},
			},
			viewKeeper:    keepers.BabylonKeeper,
			expNextCalled: true,
		},
		"custom non babylon query": {
			src: wasmvmtypes.QueryRequest{
				Custom: []byte(`{"foo":{}}`),
			},
			viewKeeper:    keepers.BabylonKeeper,
			expNextCalled: true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			var nextCalled bool
			next := QueryHandlerFn(func(ctx sdk.Context, caller sdk.AccAddress, request wasmvmtypes.QueryRequest) ([]byte, error) {
				nextCalled = true
				return nil, nil
			})

			ctx, _ := pCtx.CacheContext()
			gotData, gotErr := ChainedCustomQuerier(spec.viewKeeper, next).HandleQuery(ctx, myContractAddr, spec.src)
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, spec.expData, gotData, string(gotData))
			assert.Equal(t, spec.expNextCalled, nextCalled)
		})
	}
}

var _ viewKeeper = &MockViewKeeper{}

type MockViewKeeper struct {
	GetTestFn func(ctx sdk.Context, actor sdk.AccAddress) string
}

func (m MockViewKeeper) GetTest(ctx sdk.Context, actor sdk.AccAddress) string {
	if m.GetTestFn == nil {
		panic("not expected to be called")
	}
	return m.GetTestFn(ctx, actor)
}
