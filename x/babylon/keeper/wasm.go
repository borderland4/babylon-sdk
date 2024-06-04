package keeper

import (
	"encoding/hex"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	"github.com/babylonchain/babylon-sdk/x/babylon/contract"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SendBeginBlockMsg sends a BeginBlock sudo message to the given contract via sudo
func (k Keeper) SendBeginBlockMsg(ctx sdk.Context, contractAddr sdk.AccAddress) error {
	headerInfo := ctx.HeaderInfo()
	msg := contract.SudoMsg{
		BeginBlockMsg: &contract.BeginBlock{
			Height:     headerInfo.Height,
			HashHex:    hex.EncodeToString(headerInfo.Hash),
			Time:       headerInfo.Time,
			ChainID:    headerInfo.ChainID,
			AppHashHex: hex.EncodeToString(headerInfo.AppHash),
		},
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
