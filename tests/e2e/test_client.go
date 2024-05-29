package e2e

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm/ibctesting"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonchain/babylon-sdk/demo/app"
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
	Chain *ibctesting.TestChain
}

func NewProviderClient(t *testing.T, chain *ibctesting.TestChain) *TestProviderClient {
	return &TestProviderClient{t: t, Chain: chain}
}

func (p TestProviderClient) mustExec(contract sdk.AccAddress, payload string, funds []sdk.Coin) *abci.ExecTxResult {
	rsp, err := p.Exec(contract, payload, funds...)
	require.NoError(p.t, err)
	return rsp
}

func (p TestProviderClient) Exec(contract sdk.AccAddress, payload string, funds ...sdk.Coin) (*abci.ExecTxResult, error) {
	rsp, err := p.Chain.SendMsgs(&wasmtypes.MsgExecuteContract{
		Sender:   p.Chain.SenderAccount.GetAddress().String(),
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
	Chain     *ibctesting.TestChain
	Contracts ConsumerContract
	App       *app.ConsumerApp
}

func NewConsumerClient(t *testing.T, chain *ibctesting.TestChain) *TestConsumerClient {
	return &TestConsumerClient{t: t, Chain: chain, App: chain.App.(*app.ConsumerApp)}
}

type ConsumerContract struct {
	Babylon    sdk.AccAddress
	BTCStaking sdk.AccAddress
}

// TODO(babylon): deploy Babylon contracts
func (p *TestConsumerClient) BootstrapContracts(adminAddr sdk.AccAddress) (*ConsumerContract, error) {
	babylonContractWasmId := p.Chain.StoreCodeFile(buildPathToWasm("babylon_contract.wasm")).CodeID
	btcStakingContractWasmId := p.Chain.StoreCodeFile(buildPathToWasm("btc_staking.wasm")).CodeID

	// Instantiate the contract
	// TODO: parameterise
	initMsg := map[string]interface{}{
		"network":                         "regtest",
		"babylon_tag":                     "01020304",
		"btc_confirmation_depth":          1,
		"checkpoint_finalization_timeout": 2,
		"notify_cosmos_zone":              false,
		"btc_staking_code_id":             btcStakingContractWasmId,
		"admin":                           adminAddr.String(),
	}
	initMsgBytes, err := json.Marshal(initMsg)
	if err != nil {
		return nil, err
	}

	babylonContractAddr := InstantiateContract(p.t, p.Chain, babylonContractWasmId, initMsgBytes)
	btcStakingContractAddr := Querier(p.t, p.Chain)(babylonContractAddr.String(), Query{"config": {}})["btc_staking"]

	r := ConsumerContract{
		Babylon:    babylonContractAddr,
		BTCStaking: sdk.MustAccAddressFromBech32(btcStakingContractAddr.(string)),
	}
	p.Contracts = r
	return &r, nil
}
