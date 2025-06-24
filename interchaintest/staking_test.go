package interchaintest

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/stretchr/testify/require"

	helpers "github.com/echofi-ai/echofi/tests/interchaintest/helpers"
)

func TestStakeTokenAndUnstakeToken(t *testing.T) {
	t.Parallel()

	cfg := echofiConfig

	// Base setup
	chains := CreateChainWithCustomConfig(t, 1, 0, cfg)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	// Chains
	echofi := chains[0].(*cosmos.CosmosChain)

	// Users
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", sdkmath.NewInt(10_000_000_000), echofi, echofi)
	user := users[0]

	// vals
	vals := helpers.GetValidators(t, ctx, echofi)
	valoper := vals.Validators[0].OperatorAddress

	stakeAmt := 1_000_000
	// stake
	helpers.StakeTokens(t, ctx, echofi, user, valoper, fmt.Sprintf("%d%s", stakeAmt, echofi.Config().Denom))

	// check delegation
	delegation := helpers.QueryDelegation(t, ctx, echofi, user.FormattedAddress(), valoper)
	require.NotNil(t, delegation)
	require.Equal(t, fmt.Sprintf("%d", stakeAmt), delegation.DelegationResponse.Balance.Amount)

	// unstake
	helpers.UnstakeTokens(t, ctx, echofi, user, valoper, fmt.Sprintf("%d%s", stakeAmt, echofi.Config().Denom))

	// check unstake
	res := helpers.QueryUnbondingDelegation(t, ctx, echofi, user.FormattedAddress(), valoper)

	require.Greater(t, len(res.Unbond.Entries), 0, "no unbonding entries found")
	require.Equal(t, valoper, res.Unbond.ValidatorAddress, "validator address mismatch")

	// check delegation not found
	cmd := []string{
		"echofid", "query", "staking", "delegation", user.FormattedAddress(), valoper,
		"--output", "json",
		"--node", echofi.GetRPCAddress(),
	}

	_, _, err := echofi.Exec(ctx, cmd, nil)
	require.Error(t, err)
	require.ErrorContains(t, err, "key not found")

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
