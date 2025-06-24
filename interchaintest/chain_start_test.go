package interchaintest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicEchofiStart(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	chains := CreateThisBranchChain(t, 1, 0)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	// chain := chains[0].(*cosmos.CosmosChain)

	// userFunds := sdkmath.NewInt(10_000_000_000)
	// users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chain)
	// chainUser := users[0]

	require.NotNil(t, ic)
	require.NotNil(t, ctx)

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
