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

func (k Keeper) createNewPool(ctx sdk.Context, creatorAddress sdk.AccAddress, token1, token2 sdk.Coin, fee sdk.Dec) (pool types.Pool, err error) {

	poolId := k.GetNextPoolNumber(ctx)
	initialTotalShare, err := token1.Amount.ToLegacyDec().Mul(token2.Amount.ToLegacyDec()).ApproxSqrt()
	if err != nil {
		return pool, err
	}

	pool = types.NewPool(poolId, token1, token2, fee, creatorAddress, initialTotalShare)
	share := sdk.NewDecCoinFromDec(pool.GetPoolShareDenom(), initialTotalShare)
	poolShare := types.NewPoolShare(creatorAddress, share)

	coins := sdk.Coins{pool.Token1, pool.Token2}
	err = k.bankKeeper.SendCoins(ctx, pool.GetCreatorAddress(), pool.GetPoolAddress(), coins)
	if err != nil {
		return pool, err
	}

	k.SetPool(ctx, pool)
	k.SetPoolShare(ctx, poolShare)
	k.SetNextPoolNumber(ctx, pool.GetId()+1)
	return pool, nil
}
