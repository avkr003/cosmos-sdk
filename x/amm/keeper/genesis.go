package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/amm/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	k.SetNextPoolNumber(ctx, genState.NextPoolNumber)

	for _, pool := range genState.Pools {
		k.SetPool(ctx, pool)
	}

	for _, poolShare := range genState.PoolShares {
		k.SetPoolShare(ctx, poolShare)
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	pools := k.GetAllPools(ctx)
	poolShares := k.GetAllPoolShares(ctx)

	return &types.GenesisState{
		Params:         k.GetParams(ctx),
		NextPoolNumber: k.GetNextPoolNumber(ctx),
		Pools:          pools,
		PoolShares:     poolShares,
	}
}
