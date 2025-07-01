package interchaintest

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"testing"

	interchaintest "github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/cosmos/interchaintest/v10/testreporter"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	app "github.com/echofi-ai/echofi/app"
)

var (
	VotingPeriod     = "10s"
	MaxDepositPeriod = "10s"
	Denom            = "uecho"

	echoRepo, echoVersion = GetDockerImageInfo()

	EchofiImage = ibc.DockerImage{
		Repository: echoRepo,
		Version:    echoVersion,
		UIDGID:     "1025:1025",
	}

	// SDK v47 Genesis
	defaultGenesisKV = []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: VotingPeriod,
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: MaxDepositPeriod,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: Denom,
		},
		// {
		// 	Key:   "app_state.feepay.params.enable_feepay",
		// 	Value: false,
		// },
	}

	echofiConfig = ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "echofi",
		ChainID:             "echofi-2",
		Images:              []ibc.DockerImage{EchofiImage},
		Bin:                 "echofid",
		Bech32Prefix:        "echofi",
		Denom:               Denom,
		CoinType:            "118",
		GasPrices:           fmt.Sprintf("0%s", Denom),
		GasAdjustment:       2.0,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		ConfigFileOverrides: nil,
		EncodingConfig:      echoEncoding(),
		ModifyGenesis:       cosmos.ModifyGenesis(defaultGenesisKV),
	}

	ibcConfig = ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "ibc-chain",
		ChainID:             "ibc-1",
		Images:              []ibc.DockerImage{EchofiImage},
		Bin:                 "echofid",
		Bech32Prefix:        "echofi",
		Denom:               "uecho",
		CoinType:            "118",
		GasPrices:           fmt.Sprintf("0%s", Denom),
		GasAdjustment:       2.0,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		ConfigFileOverrides: nil,
		EncodingConfig:      echoEncoding(),
		ModifyGenesis:       cosmos.ModifyGenesis(defaultGenesisKV),
	}

	genesisWalletAmount = sdkmath.NewInt(10_000_000)
)

const (
	AccountAddressPrefix = "echofi"
)

func init() {
	// Set prefixes.
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"
	sdk.GetConfig().SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	sdk.GetConfig().SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	sdk.GetConfig().SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	sdk.GetConfig().SetCoinType(118)
}

// echoEncoding registers the Echofi specific module codecs so that the associated types and msgs
// will be supported when writing to the blocksdb sqlite database.
func echoEncoding() *testutil.TestEncodingConfig {
	tempApplication := app.NewEchofiApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		"tempDir",
		viper.New(),
		app.EVMChainID,
		app.EmptyWasmOptions,
	)
	cfg := app.MakeEncodingConfig(app.EVMChainID)

	cfg.InterfaceRegistry = tempApplication.InterfaceRegistry()
	cfg.Amino = tempApplication.LegacyAmino()
	cfg.Codec = tempApplication.AppCodec()
	cfg.TxConfig = tempApplication.GetTxConfig()
	return &cfg
}

// CreateChain generates a new chain with a custom image (useful for upgrades)
func CreateChain(t *testing.T, numVals, numFull int, img ibc.DockerImage) []ibc.Chain {
	cfg := echofiConfig
	cfg.Images = []ibc.DockerImage{img}
	return CreateChainWithCustomConfig(t, numVals, numFull, cfg)
}

// CreateThisBranchChain generates this branch's chain (ex: from the commit)
func CreateThisBranchChain(t *testing.T, numVals, numFull int) []ibc.Chain {
	return CreateChain(t, numVals, numFull, EchofiImage)
}

func CreateChainWithCustomConfig(t *testing.T, numVals, numFull int, config ibc.ChainConfig) []ibc.Chain {
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "echofi",
			ChainName:     "echofi",
			Version:       config.Images[0].Version,
			ChainConfig:   config,
			NumValidators: &numVals,
			NumFullNodes:  &numFull,
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	// chain := chains[0].(*cosmos.CosmosChain)
	return chains
}

func BuildInitialChain(t *testing.T, chains []ibc.Chain) (*interchaintest.Interchain, context.Context, *client.Client, string) {
	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain()

	for _, chain := range chains {
		ic = ic.AddChain(chain)
	}

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()
	client, network := interchaintest.DockerSetup(t)

	err := ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
		// This can be used to write to the block database which will index all block data e.g. txs, msgs, events, etc.
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
	})
	require.NoError(t, err)

	return ic, ctx, client, network
}
