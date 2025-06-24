package interchaintest

import (
	"strconv"
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/stretchr/testify/require"

	helpers "github.com/echofi-ai/echofi/tests/interchaintest/helpers"
)

func TestGov(t *testing.T) {
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

	// gov
	// ===== SUBMIT TEXT PROPOSAL =====
	propID := helpers.SubmitProposalFromFile(t, ctx, echofi, user, "/opt/proposal-test.json")
	propIDUint64, err := strconv.ParseUint(propID, 10, 64)
	require.NoError(t, err, "error converting propID to int64")

	// ===== VALIDATORS VOTES YES =====
	helpers.ValidatorVote(t, ctx, echofi, propIDUint64, 25)

	// check proposal status pass
	status := helpers.QueryProposalStatus(t, ctx, echofi, propID)
	require.Equal(t, status, "PROPOSAL_STATUS_PASSED")

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
