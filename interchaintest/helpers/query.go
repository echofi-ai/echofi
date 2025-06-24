package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/stretchr/testify/require"
	"testing"
)

func GetValidators(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain) Vals {
	var res Vals

	cmd := []string{
		"echofid", "query", "staking", "validators",
		"--output", "json",
		"--node", chain.GetRPCAddress(),
	}

	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err)

	// print stdout
	fmt.Println(string(stdout))

	// put the stdout json into res
	if err := json.Unmarshal(stdout, &res); err != nil {
		t.Fatal(err)
	}

	return res
}

func QueryDelegation(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, delegator, validator string) *DelegationQueryResponse {
	cmd := []string{
		"echofid", "query", "staking", "delegation", delegator, validator,
		"--output", "json",
		"--node", chain.GetRPCAddress(),
	}

	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err)

	// print raw response for debugging
	fmt.Println(string(stdout))

	var res DelegationQueryResponse
	if err := json.Unmarshal(stdout, &res); err != nil {
		t.Fatalf("failed to parse delegation response: %v", err)
	}

	return &res
}

func QueryUnbondingDelegation(
	t *testing.T,
	ctx context.Context,
	chain *cosmos.CosmosChain,
	delegatorAddr, validatorAddr string,
) UnbondingDelegationResponse {
	cmd := []string{
		"echofid", "query", "staking", "unbonding-delegation",
		delegatorAddr, validatorAddr,
		"--node", chain.GetRPCAddress(),
		"--output", "json",
	}

	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err, "failed to query unbonding delegation")

	fmt.Println(string(stdout))

	var res UnbondingDelegationResponse
	err = json.Unmarshal(stdout, &res)
	require.NoError(t, err, "failed to unmarshal unbonding delegation response")

	return res
}
