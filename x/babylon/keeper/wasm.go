package keeper

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	"github.com/babylonchain/babylon-sdk/x/babylon/contract"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SendTestSudoMsg sends a test sudo message to the given contract via sudo
// TODO: implement sudo messages
func (k Keeper) SendTestSudoMsg(ctx sdk.Context, contractAddr sdk.AccAddress) error {
	msg := contract.SudoMsg{
		TestSudoMsg: &struct{}{},
	}
	return k.doSudoCall(ctx, contractAddr, msg)
}

// caller must ensure gas limits are set proper and handle panics
func (k Keeper) doSudoCall(ctx sdk.Context, contractAddr sdk.AccAddress, msg contract.SudoMsg) error {
	bz, err := json.Marshal(msg)
	if err != nil {
		return errorsmod.Wrap(err, "marshal sudo msg")
	}
	_, err = k.wasm.Sudo(ctx, contractAddr, bz)
	return err
}
