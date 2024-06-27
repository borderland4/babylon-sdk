package keeper

import (
	"context"
	"encoding/hex"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	"github.com/babylonchain/babylon-sdk/x/babylon/contract"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SendBeginBlockMsg sends a BeginBlock sudo message to the BTC staking contract via sudo
func (k Keeper) SendBeginBlockMsg(c context.Context) error {
	ctx := sdk.UnwrapSDKContext(c)

	// get address of the BTC staking contract
	addrStr := k.GetParams(ctx).BtcStakingContractAddress
	if len(addrStr) == 0 {
		// the BTC staking contract address is not set yet, skip sending BeginBlockMsg
		return nil
	}
	addr := sdk.MustAccAddressFromBech32(addrStr)

	// construct the sudo message
	headerInfo := ctx.HeaderInfo()
	msg := contract.SudoMsg{
		BeginBlockMsg: &contract.BeginBlock{
			HashHex:    hex.EncodeToString(headerInfo.Hash),
			AppHashHex: hex.EncodeToString(headerInfo.AppHash),
		},
	}

	// send the sudo call
	return k.doSudoCall(ctx, addr, msg)
}

// SendEndBlockMsg sends a EndBlock sudo message to the BTC staking contract via sudo
func (k Keeper) SendEndBlockMsg(c context.Context) error {
	ctx := sdk.UnwrapSDKContext(c)

	// get address of the BTC staking contract
	addrStr := k.GetParams(ctx).BtcStakingContractAddress
	if len(addrStr) == 0 {
		// the BTC staking contract address is not set yet, skip sending EndBlockMsg
		return nil
	}
	addr := sdk.MustAccAddressFromBech32(addrStr)

	// construct the sudo message
	headerInfo := ctx.HeaderInfo()
	msg := contract.SudoMsg{
		EndBlockMsg: &contract.EndBlock{
			HashHex:    hex.EncodeToString(headerInfo.Hash),
			AppHashHex: hex.EncodeToString(headerInfo.AppHash),
		},
	}

	// send the sudo call
	return k.doSudoCall(ctx, addr, msg)
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
