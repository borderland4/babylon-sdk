package e2e

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm/ibctesting"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonchain/babylon-sdk/demo/app"
	babylon "github.com/babylonchain/babylon-sdk/x/babylon"
	"github.com/babylonchain/babylon-sdk/x/babylon/types"
)

// Query is a query type used in tests only
type Query map[string]map[string]any

// QueryResponse is a response type used in tests only
type QueryResponse map[string]any

// To can be used to navigate through the map structure
func (q QueryResponse) To(path ...string) QueryResponse {
	r, ok := q[path[0]]
	if !ok {
		panic(fmt.Sprintf("key %q does not exist", path[0]))
	}
	var x QueryResponse = r.(map[string]any)
	if len(path) == 1 {
		return x
	}
	return x.To(path[1:]...)
}

func (q QueryResponse) Array(key string) []QueryResponse {
	val, ok := q[key]
	if !ok {
		panic(fmt.Sprintf("key %q does not exist", key))
	}
	sl := val.([]any)
	result := make([]QueryResponse, len(sl))
	for i, v := range sl {
		result[i] = v.(map[string]any)
	}
	return result
}

func Querier(t *testing.T, chain *ibctesting.TestChain) func(contract string, query Query) QueryResponse {
	return func(contract string, query Query) QueryResponse {
		qRsp := make(map[string]any)
		err := chain.SmartQuery(contract, query, &qRsp)
		require.NoError(t, err)
		return qRsp
	}
}

type TestProviderClient struct {
	t     *testing.T
	chain *ibctesting.TestChain
}

func NewProviderClient(t *testing.T, chain *ibctesting.TestChain) *TestProviderClient {
	return &TestProviderClient{t: t, chain: chain}
}

func (p TestProviderClient) mustExec(contract sdk.AccAddress, payload string, funds []sdk.Coin) *sdk.Result {
	rsp, err := p.Exec(contract, payload, funds...)
	require.NoError(p.t, err)
	return rsp
}

func (p TestProviderClient) Exec(contract sdk.AccAddress, payload string, funds ...sdk.Coin) (*sdk.Result, error) {
	rsp, err := p.chain.SendMsgs(&wasmtypes.MsgExecuteContract{
		Sender:   p.chain.SenderAccount.GetAddress().String(),
		Contract: contract.String(),
		Msg:      []byte(payload),
		Funds:    funds,
	})
	return rsp, err
}

type HighLowType struct {
	High, Low int
}

// ParseHighLow convert json source type into custom type
func ParseHighLow(t *testing.T, a any) HighLowType {
	m, ok := a.(map[string]any)
	require.True(t, ok, "%T", a)
	require.Contains(t, m, "h")
	require.Contains(t, m, "l")
	h, err := strconv.Atoi(m["h"].(string))
	require.NoError(t, err)
	l, err := strconv.Atoi(m["l"].(string))
	require.NoError(t, err)
	return HighLowType{High: h, Low: l}
}

type TestConsumerClient struct {
	t         *testing.T
	chain     *ibctesting.TestChain
	contracts ConsumerContract
	app       *app.ConsumerApp
}

func NewConsumerClient(t *testing.T, chain *ibctesting.TestChain) *TestConsumerClient {
	return &TestConsumerClient{t: t, chain: chain, app: chain.App.(*app.ConsumerApp)}
}

type ConsumerContract struct {
	Babylon    sdk.AccAddress
	BTCStaking sdk.AccAddress
}

// TODO(babylon): deploy Babylon contracts
func (p *TestConsumerClient) BootstrapContracts() ConsumerContract {
	// modify end-blocker to fail fast in tests
	msModule := p.app.ModuleManager.Modules[types.ModuleName].(*babylon.AppModule)
	msModule.SetAsyncTaskRspHandler(babylon.PanicOnErrorExecutionResponseHandler())

	babylonContractWasmId := p.chain.StoreCodeFile(buildPathToWasm("babylon_contract.wasm")).CodeID
	btcStakingContractWasmId := p.chain.StoreCodeFile(buildPathToWasm("btc_staking.wasm")).CodeID

	// Instantiate the contract
	// TODO: parameterise
	initMsg := fmt.Sprintf(`{ "network": %q, "babylon_tag": %q, "btc_confirmation_depth": %d, "checkpoint_finalization_timeout": %d, "notify_cosmos_zone": %s, "btc_staking_code_id": %d }`,
		"regtest",
		"01020304",
		1,
		2,
		"false",
		btcStakingContractWasmId,
	)
	initMsgBytes := []byte(initMsg)

	babylonContractAddr := InstantiateContract(p.t, p.chain, babylonContractWasmId, initMsgBytes)
	btcStakingContractAddr := Querier(p.t, p.chain)(babylonContractAddr.String(), Query{"config": {}})["btc_staking"]

	r := ConsumerContract{
		Babylon:    babylonContractAddr,
		BTCStaking: sdk.MustAccAddressFromBech32(btcStakingContractAddr.(string)),
	}
	p.contracts = r
	return r
}
