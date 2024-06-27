package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/amm/types"
)

func (k Keeper) GetPoolShare(ctx sdk.Context, id uint64, address sdk.AccAddress) (poolShare types.PoolShare, found bool) {
	store := ctx.KVStore(k.storeKey)

	value := store.Get(types.GetPoolShareKey(id, address))
	if value == nil {
		return poolShare, false
	}

	poolShare = types.MustUnmarshalPoolShare(k.cdc, value)
	return poolShare, true
}

func (k Keeper) SetPoolShare(ctx sdk.Context, poolShare types.PoolShare) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalPoolShare(k.cdc, &poolShare)
	store.Set(poolShare.GetKey(), bz)
}

func (k Keeper) IterateAllAccountsByPool(ctx sdk.Context, id uint64, function func(poolShare types.PoolShare) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	poolSharesPrefixKey := types.GetPoolSharesKey(id)

	iterator := sdk.KVStorePrefixIterator(store, poolSharesPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		poolShare := types.MustUnmarshalPoolShare(k.cdc, iterator.Value())
		stop := function(poolShare)
		if stop {
			break
		}
	}
}

func (k Keeper) GetAllPoolSharess(ctx sdk.Context) []types.PoolShare {
	totalPools := k.GetNextPoolNumber(ctx) - 1
	if totalPools == 0 {
		return []types.PoolShare{}
	}
	store := ctx.KVStore(k.storeKey)

	poolShares := []types.PoolShare{}

	iterator := sdk.KVStorePrefixIterator(store, types.PoolShareKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		poolShare := types.MustUnmarshalPoolShare(k.cdc, iterator.Value())
		poolShares = append(poolShares, poolShare)
	}

	return poolShares
}
