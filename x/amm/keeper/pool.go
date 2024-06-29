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

	allowedTokens := k.GetParams(ctx).AllowedTokens
	token1Found := false
	token2Found := false

	for _, allowedToken := range allowedTokens {
		if token1.Denom == allowedToken {
			token1Found = true
		}
		if token2.Denom == allowedToken {
			token2Found = true
		}
		if token1Found && token2Found {
			break
		}
	}

	if !token1Found {
		return pool, types.ErrTokenNotAllowed.Wrapf(": " + token1.Denom)
	}

	if !token2Found {
		return pool, types.ErrTokenNotAllowed.Wrapf(": " + token2.Denom)
	}

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

func (k Keeper) joinPool(ctx sdk.Context, poolId uint64, fromAddress sdk.AccAddress, token sdk.Coin) (pool types.Pool, sharesAdded sdk.Dec, err error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return pool, sdk.ZeroDec(), types.ErrPoolNotFound
	}

	if pool.GetToken1().Denom != token.Denom && pool.GetToken2().Denom != token.Denom {
		return pool, sdk.ZeroDec(), types.ErrInvalidToken
	}

	if pool.GetToken1().IsZero() || pool.GetToken2().IsZero() {
		return pool, sdk.ZeroDec(), types.ErrInsufficientLiquidity.Wrapf("one or both assets of pool is empty. refill the pool or create new pool")
	}

	isRequiredToken2 := pool.GetToken1().Denom == token.Denom
	requiredCoinDenom := pool.GetToken2().Denom
	requiredRatio := pool.GetToken2().Amount.ToLegacyDec().Quo(pool.GetToken1().Amount.ToLegacyDec())
	if !isRequiredToken2 {
		requiredRatio = sdk.OneDec().Quo(requiredRatio)
		requiredCoinDenom = pool.GetToken1().Denom
	}
	requiredAmount := requiredRatio.MulInt(token.Amount).TruncateInt()
	requiredCoin := sdk.NewCoin(requiredCoinDenom, requiredAmount)

	coins := sdk.Coins{token, requiredCoin}.Sort()
	err = k.bankKeeper.SendCoins(ctx, fromAddress, pool.GetPoolAddress(), coins)
	if err != nil {
		return pool, sdk.ZeroDec(), err
	}

	newShares := requiredRatio.Mul(pool.TotalShares)
	if isRequiredToken2 {
		pool.AddToToken1(token)
		pool.AddToToken2(requiredCoin)
	} else {
		pool.AddToToken1(requiredCoin)
		pool.AddToToken2(token)
	}
	pool.AddTotalShares(newShares)
	k.SetPool(ctx, pool)

	share := sdk.NewDecCoinFromDec(pool.GetPoolShareDenom(), newShares)
	poolShare, found := k.GetPoolShare(ctx, fromAddress)
	if !found {
		poolShare = types.NewPoolShare(fromAddress, share)
	} else {
		poolShare.AddShare(share)
	}
	k.SetPoolShare(ctx, poolShare)

	return pool, newShares, nil
}

func (k Keeper) swap(ctx sdk.Context, poolId uint64, fromAddress sdk.AccAddress, tokenIn sdk.Coin) (tokenOut sdk.Coin, err error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, types.ErrPoolNotFound
	}
	if pool.GetToken1().Denom != tokenIn.Denom && pool.GetToken2().Denom != tokenIn.Denom {
		return sdk.Coin{}, types.ErrInvalidToken
	}

	if pool.GetToken1().IsZero() || pool.GetToken2().IsZero() {
		return sdk.Coin{}, types.ErrInsufficientLiquidity
	}

	effectiveTokenInAmount := tokenIn.Amount.ToLegacyDec().Mul(sdk.OneDec().Sub(pool.GetFee())).TruncateInt()

	isRequiredToken2 := pool.GetToken1().Denom == tokenIn.Denom

	if isRequiredToken2 {
		outputTokenAmount := (pool.GetToken2().Amount.Mul(effectiveTokenInAmount)).Quo(pool.GetToken1().Amount.Add(effectiveTokenInAmount))
		tokenOut = sdk.NewCoin(pool.GetToken2().Denom, outputTokenAmount)
		pool.AddToToken1(tokenIn)
		err = pool.SubtractFromToken2(tokenOut)
		if err != nil {
			return sdk.Coin{}, types.ErrInvalidToken
		}
	} else {
		outputTokenAmount := (pool.GetToken1().Amount.Mul(effectiveTokenInAmount)).Quo(pool.GetToken2().Amount.Add(effectiveTokenInAmount))
		tokenOut = sdk.NewCoin(pool.GetToken1().Denom, outputTokenAmount)
		pool.AddToToken2(tokenIn)
		err = pool.SubtractFromToken1(tokenOut)
		if err != nil {
			return sdk.Coin{}, types.ErrInvalidToken
		}
	}

	err = k.bankKeeper.SendCoins(ctx, fromAddress, pool.GetPoolAddress(), sdk.NewCoins(tokenIn))
	if err != nil {
		return sdk.Coin{}, err
	}
	err = k.bankKeeper.SendCoins(ctx, pool.GetPoolAddress(), fromAddress, sdk.NewCoins(tokenOut))
	if err != nil {
		return sdk.Coin{}, err
	}

	k.SetPool(ctx, pool)

	return tokenOut, nil
}

func (k Keeper) exitPool(ctx sdk.Context, poolId uint64, fromAddress sdk.AccAddress, lpShares sdk.Dec, withdrawAll bool) (err error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return types.ErrPoolNotFound
	}

	poolShare, found := k.GetPoolShare(ctx, fromAddress)
	if !found {
		return types.ErrLpSharesNotFound
	}

	totalAddressShares, foundAt := poolShare.GetShare(pool.GetPoolShareDenom())
	if foundAt == -1 {
		return types.ErrLpSharesNotFound
	}

	refundShareAmount := lpShares

	if refundShareAmount.GT(totalAddressShares.Amount) {
		return types.ErrRedeemingMoreThanAllowed
	}

	if withdrawAll {
		refundShareAmount = totalAddressShares.Amount
	}

	ratio := refundShareAmount.Quo(pool.GetTotalShares()) // divided by 0 not possible here because ErrLpSharesNotFound will happen beforehand
	token1Out := sdk.NewCoin(pool.GetToken1().Denom, ratio.Mul(pool.GetToken1().Amount.ToLegacyDec()).TruncateInt())
	token2Out := sdk.NewCoin(pool.GetToken2().Denom, ratio.Mul(pool.GetToken2().Amount.ToLegacyDec()).TruncateInt())

	err = pool.SubtractFromToken1(token1Out)
	if err != nil {
		return err
	}
	err = pool.SubtractFromToken2(token2Out)
	if err != nil {
		return err
	}

	tokensOut := sdk.NewCoins(token1Out, token2Out)
	err = k.bankKeeper.SendCoins(ctx, pool.GetPoolAddress(), fromAddress, tokensOut)
	if err != nil {
		return err
	}

	err = poolShare.SubtractShare(sdk.NewDecCoinFromDec(pool.GetPoolShareDenom(), refundShareAmount))
	if err != nil {
		return err
	}

	pool.SubtractFromTotalShares(refundShareAmount)

	k.SetPool(ctx, pool)
	k.SetPoolShare(ctx, poolShare)
	return nil
}

func (k Keeper) refillEmptyPool(ctx sdk.Context, fromAddress sdk.AccAddress, poolId uint64, token1, token2 sdk.Coin) (pool types.Pool, err error) {

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return pool, types.ErrPoolNotFound
	}

	if !(pool.GetToken1().IsZero() && pool.GetToken2().IsZero()) {
		return pool, types.ErrPoolNotEmpty
	}

	totalShare, err := token1.Amount.ToLegacyDec().Mul(token2.Amount.ToLegacyDec()).ApproxSqrt()
	if err != nil {
		return pool, err
	}

	pool = types.NewPool(poolId, token1, token2, pool.GetFee(), pool.GetCreatorAddress(), totalShare)
	share := sdk.NewDecCoinFromDec(pool.GetPoolShareDenom(), totalShare)
	poolShare := types.NewPoolShare(pool.GetCreatorAddress(), share)

	coins := sdk.Coins{pool.Token1, pool.Token2}
	err = k.bankKeeper.SendCoins(ctx, fromAddress, pool.GetPoolAddress(), coins)
	if err != nil {
		return pool, err
	}

	k.SetPool(ctx, pool)
	k.SetPoolShare(ctx, poolShare)
	return pool, nil
}
