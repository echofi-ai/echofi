package helpers

import (
	"context"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/cosmos/interchaintest/v10/testutil"
	"github.com/stretchr/testify/require"
)

func StakeTokens(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, user ibc.Wallet, valoper, coinAmt string) {
	// amount is #utoken
	cmd := []string{
		"echofid", "tx", "staking", "delegate", valoper, coinAmt,
		"--home", chain.HomeDir(),
		"--from", user.KeyName(),
		"--node", chain.GetRPCAddress(),
		"--chain-id", chain.Config().ChainID,
		"--gas", "500000",
		"--keyring-dir", chain.HomeDir(),
		"--keyring-backend", keyring.BackendTest,
		"-y",
	}
	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	if err := testutil.WaitForBlocks(ctx, 2, chain); err != nil {
		t.Fatal(err)
	}
}

func UnstakeTokens(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, user ibc.Wallet, valoper, coinAmt string) {
	cmd := []string{
		"echofid", "tx", "staking", "unbond", valoper, coinAmt,
		"--home", chain.HomeDir(),
		"--from", user.KeyName(),
		"--node", chain.GetRPCAddress(),
		"--chain-id", chain.Config().ChainID,
		"--gas", "500000",
		"--keyring-dir", chain.HomeDir(),
		"--keyring-backend", keyring.BackendTest,
		"-y",
	}

	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	if err := testutil.WaitForBlocks(ctx, 2, chain); err != nil {
		t.Fatal(err)
	}
}

func ClaimStakingRewards(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, user ibc.Wallet, valoper string) {
	cmd := []string{
		"echofid", "tx", "distribution", "withdraw-rewards", valoper,
		"--home", chain.HomeDir(),
		"--from", user.KeyName(),
		"--node", chain.GetRPCAddress(),
		"--chain-id", chain.Config().ChainID,
		"--gas", "500000",
		"--keyring-dir", chain.HomeDir(),
		"--keyring-backend", keyring.BackendTest,
		"-y",
	}
	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	if err := testutil.WaitForBlocks(ctx, 2, chain); err != nil {
		t.Fatal(err)
	}
}
