package keeper

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/babylonchain/babylon-sdk/x/babylon/contract"
	"github.com/babylonchain/babylon-sdk/x/babylon/types"
)

// AuthSource abstract type that provides contract authorization.
// This is an extension point for custom implementations.
type AuthSource interface {
	// IsAuthorized returns if the contract authorized to execute a virtual stake message
	IsAuthorized(ctx sdk.Context, contractAddr sdk.AccAddress) bool
}

// abstract keeper
type msKeeper interface{}

type CustomMsgHandler struct {
	k    msKeeper
	auth AuthSource
}

// NewDefaultCustomMsgHandler constructor to set up the CustomMsgHandler with default max cap authorization
func NewDefaultCustomMsgHandler(k *Keeper) *CustomMsgHandler {
	return &CustomMsgHandler{k: k, auth: defaultMaxCapAuthorizator(k)}
}

// NewCustomMsgHandler constructor to set up CustomMsgHandler with an individual auth source.
// This is an extension point for non default contract authorization logic.
func NewCustomMsgHandler(k msKeeper, auth AuthSource) *CustomMsgHandler {
	return &CustomMsgHandler{k: k, auth: auth}
}

func defaultMaxCapAuthorizator(k *Keeper) AuthSourceFn {
	return func(ctx sdk.Context, contractAddr sdk.AccAddress) bool {
		return true
	}
}

// DispatchMsg handle contract message of type Custom in the babylon namespace
func (h CustomMsgHandler) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, _ string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom == nil {
		return nil, nil, wasmtypes.ErrUnknownMsg
	}
	var customMsg contract.CustomMsg
	if err := json.Unmarshal(msg.Custom, &customMsg); err != nil {
		return nil, nil, sdkerrors.ErrJSONUnmarshal.Wrap("custom message")
	}
	if customMsg.Test == nil {
		// not our message type
		return nil, nil, wasmtypes.ErrUnknownMsg
	}

	if !h.auth.IsAuthorized(ctx, contractAddr) {
		return nil, nil, sdkerrors.ErrUnauthorized.Wrapf("contract has no permission for Babylon operations")
	}

	return h.handleTestMsg(ctx, contractAddr, customMsg.Test)
}

func (h CustomMsgHandler) handleTestMsg(ctx sdk.Context, actor sdk.AccAddress, testMsg *contract.TestMsg) ([]sdk.Event, [][]byte, error) {
	return []sdk.Event{}, nil, nil
}

// AuthSourceFn is helper for simple AuthSource types
type AuthSourceFn func(ctx sdk.Context, contractAddr sdk.AccAddress) bool

// IsAuthorized returns if the contract authorized to execute a virtual stake message
func (a AuthSourceFn) IsAuthorized(ctx sdk.Context, contractAddr sdk.AccAddress) bool {
	return a(ctx, contractAddr)
}

// abstract keeper
type integrityHandlerSource interface {
	CanInvokeStakingMsg(ctx sdk.Context, actor sdk.AccAddress) bool
}

// NewIntegrityHandler prevents any contract with max cap set to use staking
// or stargate messages. This ensures that staked "virtual" tokens are not bypassing
// the instant undelegate and burn mechanism provided by babylon.
//
// This handler should be chained before any other.
func NewIntegrityHandler(k integrityHandlerSource) wasmkeeper.MessageHandlerFunc {
	return func(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (
		events []sdk.Event,
		data [][]byte,
		err error,
	) {
		if msg.Stargate == nil && msg.Staking == nil ||
			!k.CanInvokeStakingMsg(ctx, contractAddr) {
			return nil, nil, wasmtypes.ErrUnknownMsg // pass down the chain
		}
		// reject
		return nil, nil, types.ErrUnsupported.Wrap("message type for contracts with max cap set")
	}
}

func (k Keeper) CanInvokeStakingMsg(ctx sdk.Context, actor sdk.AccAddress) bool {
	// TODO: implement access control
	return true
}
