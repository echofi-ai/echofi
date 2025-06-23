package app

import (
	"encoding/json"
	"sort"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	storetypes "cosmossdk.io/store/types"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ExportAppStateAndValidators exports the state of the application for a genesis
// file.
func (app *EchofiApp) ExportAppStateAndValidators(
	forZeroHeight bool,
	jailAllowedAddrs []string,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	// as if they could withdraw from the start of the next block
	ctx := app.BaseApp.NewContextLegacy(true, tmproto.Header{Height: app.BaseApp.LastBlockHeight()})

	// We export at last height + 1, because that's the height at which
	// Tendermint will start InitChain.
	height := app.BaseApp.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	}

	genState, err := app.mm.ExportGenesisForModules(ctx, app.appCodec, modulesToExport)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(ctx, app.AppKeepers.StakingKeeper)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	sort.SliceStable(validators, func(i, j int) bool {
		return validators[i].Power > validators[j].Power
	})
	// we have to trim this to only active consensus validators
	maxVals := app.AppKeepers.ProviderKeeper.GetMaxProviderConsensusValidators(ctx)
	if len(validators) > int(maxVals) {
		validators = validators[:maxVals]
	}

	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, err
}

// prepare for fresh start at zero height
// NOTE zero height genesis is a temporary feature which will be deprecated
// in favour of export at a block height
func (app *EchofiApp) prepForZeroHeightGenesis(ctx sdk.Context, jailAllowedAddrs []string) {
	// check if there is an allowed address list
	applyAllowedAddrs := len(jailAllowedAddrs) > 0

	allowedAddrsMap := make(map[string]bool)

	for _, addr := range jailAllowedAddrs {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		allowedAddrsMap[addr] = true
	}

	/* Handle fee distribution state. */

	// withdraw all validator commission
	err := app.AppKeepers.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		valAddr, err := app.AppKeepers.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			app.BaseApp.Logger().Error(err.Error(), "ValOperatorAddress", val.GetOperator())
		}
		_, err = app.AppKeepers.DistrKeeper.WithdrawValidatorCommission(ctx, valAddr)
		if err != nil {
			app.BaseApp.Logger().Error(err.Error(), "ValOperatorAddress", val.GetOperator())
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	// withdraw all delegator rewards
	dels, err := app.AppKeepers.StakingKeeper.GetAllDelegations(ctx)
	if err != nil {
		panic(err)
	}
	for _, delegation := range dels {
		valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		delAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		_, err = app.AppKeepers.DistrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
		if err != nil {
			panic(err)
		}
	}

	// clear validator slash events
	app.AppKeepers.DistrKeeper.DeleteAllValidatorSlashEvents(ctx)

	// clear validator historical rewards
	app.AppKeepers.DistrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reinitialize all validators (v0.46 version)
	// app.AppKeepers.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
	// 	// donate any unwithdrawn outstanding reward fraction tokens to the community pool
	// 	scraps := app.AppKeepers.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, val.GetOperator())
	// 	feePool := app.AppKeepers.DistrKeeper.GetFeePool(ctx)
	// 	feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
	// 	app.AppKeepers.DistrKeeper.SetFeePool(ctx, feePool)

	// 	err := app.AppKeepers.DistrKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return false
	// })

	// reinitialize all validators
	err = app.AppKeepers.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		// donate any unwithdrawn outstanding reward fraction tokens to the community pool
		valAddr, err := app.AppKeepers.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			panic(err)
		}
		scraps, err := app.AppKeepers.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr)
		if err != nil {
			panic(err)
		}
		feePool, err := app.AppKeepers.DistrKeeper.FeePool.Get(ctx)
		if err != nil {
			panic(err)
		}
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		err = app.AppKeepers.DistrKeeper.FeePool.Set(ctx, feePool)
		if err != nil {
			panic(err)
		}
		if err := app.AppKeepers.DistrKeeper.Hooks().AfterValidatorCreated(ctx, valAddr); err != nil {
			panic(err)
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	// reinitialize all delegations
	for _, del := range dels {
		valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delAddr, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		if err := app.AppKeepers.DistrKeeper.Hooks().BeforeDelegationCreated(ctx, delAddr, valAddr); err != nil {
			panic(err)
		}
		if err := app.AppKeepers.DistrKeeper.Hooks().AfterDelegationModified(ctx, delAddr, valAddr); err != nil {
			panic(err)
		}
	}

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */

	// iterate through redelegations, reset creation height
	err = app.AppKeepers.StakingKeeper.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		if err := app.AppKeepers.StakingKeeper.SetRedelegation(ctx, red); err != nil {
			panic(err)
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	// iterate through unbonding delegations, reset creation height
	err = app.AppKeepers.StakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		if err := app.AppKeepers.StakingKeeper.SetUnbondingDelegation(ctx, ubd); err != nil {
			panic(err)
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(app.AppKeepers.GetKey(stakingtypes.StoreKey))
	iter := storetypes.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)

	counter := int16(0)

	// Closure to ensure iterator doesn't leak.
	func() {
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			addr := sdk.ValAddress(stakingtypes.AddressFromValidatorsKey(iter.Key()))
			validator, err := app.AppKeepers.StakingKeeper.GetValidator(ctx, addr)
			if err != nil {
				panic("expected validator, not found")
			}

			validator.UnbondingHeight = 0
			if applyAllowedAddrs && !allowedAddrsMap[addr.String()] {
				validator.Jailed = true
			}

			if err = app.AppKeepers.StakingKeeper.SetValidator(ctx, validator); err != nil {
				panic(err)
			}

			counter++
		}
	}()

	_, err = app.AppKeepers.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	/* Handle slashing state. */

	// reset start height on signing infos
	err = app.AppKeepers.SlashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			if err = app.AppKeepers.SlashingKeeper.SetValidatorSigningInfo(ctx, addr, info); err != nil {
				panic(err)
			}
			return false
		},
	)
	if err != nil {
		panic(err)
	}
}
