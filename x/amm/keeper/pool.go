package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/amm/types"
	gogotypes "github.com/cosmos/gogoproto/types"
)

func (k Keeper) GetPool(ctx sdk.Context, id uint64) (pool types.Pool, found bool) {
	store := ctx.KVStore(k.storeKey)

	value := store.Get(types.GetPoolKey(id))
	if value == nil {
		return pool, false
	}

	pool = types.MustUnmarshalPool(k.cdc, value)
	return pool, true
}

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalPool(k.cdc, &pool)
	store.Set(types.GetPoolKey(pool.GetId()), bz)
}

func (k Keeper) GetAllPools(ctx sdk.Context) []types.Pool {
	totalPools := k.GetNextPoolNumber(ctx) - 1
	if totalPools == 0 {
		return []types.Pool{}
	}
	store := ctx.KVStore(k.storeKey)

	pools := make([]types.Pool, totalPools)

	iterator := sdk.KVStorePrefixIterator(store, types.PoolKey)
	defer iterator.Close()

	for i := 0; iterator.Valid(); iterator.Next() {
		pool := types.MustUnmarshalPool(k.cdc, iterator.Value())
		pools[i] = pool
		i++
	}

	return pools
}

func (k Keeper) SetNextPoolNumber(ctx sdk.Context, poolNumber uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: poolNumber})
	store.Set(types.PoolNumberKey, bz)
}

func (k Keeper) GetNextPoolNumber(ctx sdk.Context) uint64 {
	var poolNumber uint64
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.PoolNumberKey)
	if bz == nil {
		panic(fmt.Errorf("pool has not been initialized -- Should have been done in InitGenesis"))
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}

		poolNumber = val.GetValue()
	}

	return poolNumber
}

func (k Keeper) createNewPool(ctx sdk.Context, pool types.Pool, poolShare types.PoolShare) error {

	coins := sdk.Coins{pool.Token_1, pool.Token_2}
	err := k.bankKeeper.SendCoins(ctx, pool.GetCreatorAddress(), pool.GetPoolAddress(), coins)
	if err != nil {
		return err
	}

	k.SetPool(ctx, pool)
	k.SetPoolShare(ctx, poolShare)
	k.SetNextPoolNumber(ctx, pool.GetId()+1)
	return nil
}
