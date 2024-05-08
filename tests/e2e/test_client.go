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
	staking   sdk.AccAddress
	priceFeed sdk.AccAddress
	converter sdk.AccAddress
}

// TODO(babylon): deploy Babylon contracts
func (p *TestConsumerClient) BootstrapContracts() ConsumerContract {
	// modify end-blocker to fail fast in tests
	msModule := p.app.ModuleManager.Modules[types.ModuleName].(*babylon.AppModule)
	msModule.SetAsyncTaskRspHandler(babylon.PanicOnErrorExecutionResponseHandler())

	var ( // todo: configure
		tokenRatio  = "0.5"
		discount    = "0.1"
		remoteDenom = sdk.DefaultBondDenom
	)
	codeID := p.chain.StoreCodeFile(buildPathToWasm("mesh_simple_price_feed.wasm")).CodeID
	initMsg := []byte(fmt.Sprintf(`{"native_per_foreign": "%s"}`, tokenRatio))
	priceFeedContract := InstantiateContract(p.t, p.chain, codeID, initMsg)
	// virtual staking is setup by the consumer
	virtStakeCodeID := p.chain.StoreCodeFile(buildPathToWasm("mesh_virtual_staking.wasm")).CodeID
	// instantiate converter
	codeID = p.chain.StoreCodeFile(buildPathToWasm("mesh_converter.wasm")).CodeID
	initMsg = []byte(fmt.Sprintf(`{"price_feed": %q, "discount": %q, "remote_denom": %q,"virtual_staking_code_id": %d}`,
		priceFeedContract.String(), discount, remoteDenom, virtStakeCodeID))
	converterContract := InstantiateContract(p.t, p.chain, codeID, initMsg)

	staking := Querier(p.t, p.chain)(converterContract.String(), Query{"config": {}})["virtual_staking"]
	r := ConsumerContract{
		staking:   sdk.MustAccAddressFromBech32(staking.(string)),
		priceFeed: priceFeedContract,
		converter: converterContract,
	}
	p.contracts = r
	return r
}
