package babylon

import (
	"time"

	"github.com/babylonchain/babylon-sdk/x/babylon/keeper"
	"github.com/babylonchain/babylon-sdk/x/babylon/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k *keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	babylonContractAddr := sdk.MustAccAddressFromBech32(k.GetParams(ctx).BabylonContractAddress)
	if babylonContractAddr.String() == types.EmptyAddr {
		// the Babylon contract address is not set yet, skip sending BeginBlockMsg
		return nil
	}
	return k.SendBeginBlockMsg(ctx, babylonContractAddr)
}

// EndBlocker is called after every block
func EndBlocker(ctx sdk.Context, k *keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
}
