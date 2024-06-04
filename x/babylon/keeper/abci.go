package keeper

import (
	"time"

	"github.com/babylonchain/babylon-sdk/x/babylon/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) BeginBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	addrStr := k.GetParams(ctx).BabylonContractAddress
	if len(addrStr) == 0 {
		// the Babylon contract address is not set yet, skip sending BeginBlockMsg
		return nil
	}
	babylonContractAddr := sdk.MustAccAddressFromBech32(addrStr)
	return k.SendBeginBlockMsg(ctx, babylonContractAddr)
}

// EndBlocker is called after every block
func (k *Keeper) EndBlocker(ctx sdk.Context) ([]abci.ValidatorUpdate, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	return []abci.ValidatorUpdate{}, nil
}
