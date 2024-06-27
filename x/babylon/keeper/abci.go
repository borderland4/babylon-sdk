package keeper

import (
	"fmt"
	rt "runtime"
	"strings"
	"time"

	"github.com/babylonchain/babylon-sdk/x/babylon/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func printStackTrace() {
	var pcs [32]uintptr
	n := rt.Callers(2, pcs[:])
	frames := rt.CallersFrames(pcs[:n])

	var sb strings.Builder
	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	fmt.Println(sb.String())
}

func (k *Keeper) BeginBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	return k.SendBeginBlockMsg(ctx)
}

// EndBlocker is called after every block
func (k *Keeper) EndBlocker(ctx sdk.Context) ([]abci.ValidatorUpdate, error) {
	fmt.Println("Keeper EndBlocker is called") // Basic print statement
	printStackTrace()                          // Print the stack trace
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Check if the logger is set up correctly
	if k.Logger(ctx) == nil {
		fmt.Println("Logger is not set")
	} else {
		k.Logger(ctx).Info("Debug: EndBlocker called", "height", ctx.BlockHeight())
		fmt.Printf("Debug: EndBlocker called, height: %d\n", ctx.BlockHeight())
	}

	k.Logger(ctx).Info("Debug: EndBlocker called", "height", ctx.BlockHeight())
	if err := k.SendEndBlockMsg(ctx); err != nil {
		return []abci.ValidatorUpdate{}, err
	}

	return []abci.ValidatorUpdate{}, nil
}
