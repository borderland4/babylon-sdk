package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMVP(t *testing.T) {
	x := setupExampleChains(t)
	consumerCli, consumerContracts, providerCli := setupBabylonIntegration(t, x)
	require.NotEmpty(t, consumerCli.chain.ChainID)
	require.NotEmpty(t, providerCli.chain.ChainID)
	require.False(t, consumerContracts.Babylon.Empty())
	require.False(t, consumerContracts.BTCStaking.Empty())
}
