package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/cosmos/interchaintest/v10/testutil"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

// QueryVotes queries votes for a given proposal ID
func QueryVotes(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, proposalID string) []Vote {
	cmd := []string{
		"echofid", "query", "gov", "votes", proposalID,
		"--node", chain.GetRPCAddress(),
		"--output", "json",
	}

	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err, "failed to query gov votes")

	var result QueryVotesResult
	if err := json.Unmarshal(stdout, &result); err != nil {
		t.Fatalf("failed to unmarshal votes: %v", err)
	}

	// Debug output
	fmt.Printf("Votes for proposal %s:\n%s\n", proposalID, string(stdout))

	return result.Votes
}

func QueryProposalStatus(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, proposalID string) string {
	cmd := []string{
		"echofid", "query", "gov", "proposal", proposalID,
		"--node", chain.GetRPCAddress(),
		"--output", "json",
	}

	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err, "failed to query proposal")

	var wrapper struct {
		Proposal Proposal `json:"proposal"`
	}

	err = json.Unmarshal(stdout, &wrapper)
	require.NoError(t, err, "failed to unmarshal proposal")

	fmt.Printf("Proposal %s status: %s\n", proposalID, wrapper.Proposal.Status)
	return wrapper.Proposal.Status
}

// Modified from ictest
func VoteOnProposalAllValidators(ctx context.Context, c *cosmos.CosmosChain, proposalID uint64, vote string) error {
	var eg errgroup.Group
	valKey := "validator"
	for _, n := range c.Nodes() {
		if n.Validator {
			n := n
			eg.Go(func() error {
				// gas-adjustment was using 1.3 default instead of the setup's 2.0+ for some reason.
				// return n.VoteOnProposal(ctx, valKey, proposalID, vote)

				_, err := n.ExecTx(ctx, valKey,
					"gov", "vote",
					strconv.Itoa(int(proposalID)), vote, "--gas", "auto", "--gas-adjustment", "2.0",
				)
				return err
			})
		}
	}
	return eg.Wait()
}

func ValidatorVote(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, proposalID uint64, searchHeightDelta int64) {
	err := VoteOnProposalAllValidators(ctx, chain, proposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to vote on proposal")

	height, err := chain.Height(ctx)
	require.NoError(t, err, "failed to get height")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+searchHeightDelta, proposalID, govtypes.StatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")
}

func SubmitProposalFromFile(
	t *testing.T,
	ctx context.Context,
	chain *cosmos.CosmosChain,
	user ibc.Wallet,
	filePath string,
) string {
	cmd := []string{
		"echofid", "tx", "gov", "submit-proposal", filePath,
		"--keyring-backend", keyring.BackendTest,
		"--from", user.KeyName(),
		"--home", chain.HomeDir(),
		"--node", chain.GetRPCAddress(),
		"--chain-id", chain.Config().ChainID,
		"--gas", "500000",
		"--fees", "5000" + chain.Config().Denom,
		"-y",
		"--output", "json",
	}

	stdout, stderr, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err, "failed submitting proposal:\nstdout: %s\nstderr: %s", stdout, stderr)

	// Optional: log output
	fmt.Println("Submit proposal stdout:", string(stdout))

	var res struct {
		TxHash string `json:"txhash"`
	}

	err = json.Unmarshal(stdout, &res)
	require.NoError(t, err)
	txHash := res.TxHash

	err = testutil.WaitForBlocks(ctx, 2, chain)
	require.NoError(t, err)

	return GetProposalIDFromTx(t, ctx, chain, txHash)
}

func GetProposalIDFromTx(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, txHash string) string {
	cmd := []string{
		"echofid", "query", "tx", txHash,
		"--node", chain.GetRPCAddress(),
		"--output", "json",
	}
	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err)

	fmt.Println("Get tx hash stdout:", string(stdout))

	return ExtractProposalIDFromTxResponse(t, stdout)
}

func ExtractProposalIDFromTxResponse(t *testing.T, txResp []byte) string {
	var parsed struct {
		Logs   []interface{} `json:"logs"`
		Events []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"attributes"`
		} `json:"events"`
	}

	require.NoError(t, json.Unmarshal(txResp, &parsed))

	for _, ev := range parsed.Events {
		if ev.Type == "submit_proposal" {
			for _, attr := range ev.Attributes {
				if attr.Key == "proposal_id" {
					return attr.Value
				}
			}
		}
	}

	t.Fatal("proposal_id not found in tx response")
	return ""
}
