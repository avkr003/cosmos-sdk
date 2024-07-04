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

func (k Keeper) mintLpShares(ctx sdk.Context, amount sdk.Coin, toAddress sdk.AccAddress) error {
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddress, sdk.NewCoins(amount))
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) burnLpShares(ctx sdk.Context, amount sdk.Coin, fromAddress sdk.AccAddress) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAddress, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) createNewPool(ctx sdk.Context, creatorAddress sdk.AccAddress, tokens sdk.Coins, fee sdk.Dec) (pool types.Pool, err error) {

	params := k.GetParams(ctx)
	// tokens are always sorted due to validation check, even when creating pools so no need to check denom name for token 1 and token 2

	token1 := tokens[0]
	token1Found := false
	token2 := tokens[1]
	token2Found := false

	for _, allowedToken := range params.SwapAllowedTokens {
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

	// no need to check for fee to be less than 0, already checked in ValidateBasic of the message
	if fee.GT(params.MaxSwapFee) {
		return pool, types.ErrInvalidFees.Wrapf("swap fee is greater than allowed %s ", params.MaxSwapFee.String())
	}

	poolId := k.GetNextPoolNumber(ctx)

	initialTotalShare, err := types.GetInitialPoolShares(token1, token2)
	if err != nil {
		return pool, err
	}

	pool = types.NewPool(poolId, token1, token2, fee, creatorAddress)

	share := sdk.NewCoin(pool.GetPoolShareDenom(), initialTotalShare)
	err = k.mintLpShares(ctx, share, creatorAddress)
	if err != nil {
		return types.Pool{}, err
	}

	coins := sdk.Coins{pool.Token1, pool.Token2}
	err = k.bankKeeper.SendCoins(ctx, pool.GetCreatorAddress(), pool.GetPoolAddress(), coins)
	if err != nil {
		return pool, err
	}

	k.SetPool(ctx, pool)
	k.SetNextPoolNumber(ctx, pool.GetId()+1)
	return pool, nil
}

func (k Keeper) joinPool(ctx sdk.Context, poolId uint64, fromAddress sdk.AccAddress, tokens sdk.Coins) (pool types.Pool, err error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return pool, types.ErrPoolNotFound
	}

	if len(tokens) > 2 || len(tokens) == 0 {
		return pool, types.ErrInvalidTokens.Wrapf("only 1 or 2 token can be given")
	}

	if len(tokens) == 2 {
		// tokens are always sorted due to validation check, even when creating pools so no need to check denom name for token 1 and token 2
		if tokens[0].Denom != pool.GetToken1().Denom || tokens[1].Denom != pool.GetToken2().Denom {
			return pool, types.ErrInvalidTokens.Wrapf("tokens do not match pool tokens")
		}
		if pool.GetToken1().IsZero() && pool.GetToken2().IsZero() {
			err = k.refillEmptyPool(ctx, fromAddress, pool, tokens)
			return pool, err
		}
		token1 := tokens[0]
		token2 := tokens[1]
		current2by1Ratio := pool.GetToken2by1Ratio()
		expectedToken2 := token1.Amount.ToLegacyDec().Mul(current2by1Ratio).RoundInt()

		// token 2 amount given is greater than expected token 2 required for given token 1
		if token2.Amount.GT(expectedToken2) {
			err = k.singleTokenJoinPool(ctx, fromAddress, pool, token1)
			return pool, err
		} else {
			// token 1 amount given is greater than expected token 1 required for given token 2, or it's given in required ratio - then it doesn't matter which token is send to function
			err = k.singleTokenJoinPool(ctx, fromAddress, pool, token2)
			return pool, err
		}
	}

	if len(tokens) == 1 {
		err = k.singleTokenJoinPool(ctx, fromAddress, pool, tokens[0])
		return pool, err
	}
	return pool, err
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

	effectiveTokenInAmount := tokenIn.Amount.ToLegacyDec().Mul(sdk.OneDec().Sub(pool.GetFee())).RoundInt()

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

func (k Keeper) exitPool(ctx sdk.Context, poolId uint64, fromAddress sdk.AccAddress, lpShares sdk.Int, withdrawAll bool) (err error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return types.ErrPoolNotFound
	}

	totalAccountShares := k.bankKeeper.GetBalance(ctx, fromAddress, pool.GetPoolShareDenom())
	if totalAccountShares.IsZero() {
		return types.ErrLpSharesNotFound
	}

	refundShareAmount := lpShares

	if withdrawAll {
		refundShareAmount = totalAccountShares.Amount
	}

	if refundShareAmount.GT(totalAccountShares.Amount) {
		return types.ErrRedeemingMoreThanAllowed
	}

	totalShares := k.bankKeeper.GetSupply(ctx, pool.GetPoolShareDenom())
	if totalShares.Amount.IsZero() {
		return types.ErrEmptyPool
	}
	sharesRatio := refundShareAmount.ToLegacyDec().Quo(totalShares.Amount.ToLegacyDec()) // divided by 0 not possible here because ErrLpSharesNotFound will happen beforehand
	token1Out := sdk.NewCoin(pool.GetToken1().Denom, sharesRatio.MulInt(pool.GetToken1().Amount).RoundInt())
	token2Out := sdk.NewCoin(pool.GetToken2().Denom, sharesRatio.MulInt(pool.GetToken2().Amount).RoundInt())

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

	k.SetPool(ctx, pool)

	err = k.burnLpShares(ctx, sdk.NewCoin(pool.GetPoolShareDenom(), refundShareAmount), fromAddress)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) singleTokenJoinPool(ctx sdk.Context, fromAddress sdk.AccAddress, pool types.Pool, token sdk.Coin) error {

	if pool.GetToken1().Denom != token.Denom && pool.GetToken2().Denom != token.Denom {
		return types.ErrInvalidToken
	}

	if pool.GetToken1().IsZero() || pool.GetToken2().IsZero() {
		return types.ErrInsufficientLiquidity.Wrapf("one or both assets of pool is empty. refill the pool or create new pool")
	}

	isRequiredToken2 := pool.GetToken1().Denom == token.Denom
	requiredCoinDenom := pool.GetToken2().Denom
	requiredTokenRatio := pool.GetToken2by1Ratio()
	requiredShareRatio := token.Amount.ToLegacyDec().Quo(pool.GetToken1().Amount.ToLegacyDec())
	if !isRequiredToken2 {
		requiredTokenRatio = sdk.OneDec().Quo(requiredTokenRatio)
		requiredCoinDenom = pool.GetToken1().Denom
		requiredShareRatio = token.Amount.ToLegacyDec().Quo(pool.GetToken2().Amount.ToLegacyDec())
	}
	requiredTokenAmount := requiredTokenRatio.MulInt(token.Amount).RoundInt()
	requiredCoin := sdk.NewCoin(requiredCoinDenom, requiredTokenAmount)

	coins := sdk.Coins{token, requiredCoin}.Sort()
	err := k.bankKeeper.SendCoins(ctx, fromAddress, pool.GetPoolAddress(), coins)
	if err != nil {
		return err
	}

	totalShares := k.bankKeeper.GetSupply(ctx, pool.GetPoolShareDenom())
	if totalShares.Amount.IsZero() {
		return types.ErrEmptyPool
	}

	newShares := requiredShareRatio.MulInt(totalShares.Amount).RoundInt()
	if isRequiredToken2 {
		pool.AddToToken1(token)
		pool.AddToToken2(requiredCoin)
	} else {
		pool.AddToToken1(requiredCoin)
		pool.AddToToken2(token)
	}
	k.SetPool(ctx, pool)

	share := sdk.NewCoin(pool.GetPoolShareDenom(), newShares)
	err = k.mintLpShares(ctx, share, fromAddress)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) refillEmptyPool(ctx sdk.Context, fromAddress sdk.AccAddress, pool types.Pool, tokens sdk.Coins) error {

	token1 := tokens[0]
	token2 := tokens[1]

	poolToken1Denom, poolToken2Denom := pool.GetPoolTokensFromName()
	if token1.Denom != poolToken1Denom || token2.Denom != poolToken2Denom {
		return types.ErrInvalidToken.Wrapf("token denom does not match pool name tokens")
	}

	initialTotalShare, err := types.GetInitialPoolShares(token1, token2)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoins(ctx, fromAddress, pool.GetPoolAddress(), sdk.NewCoins(token1, token2))
	if err != nil {
		return err
	}

	pool = types.NewPool(pool.GetId(), token1, token2, pool.GetFee(), pool.GetCreatorAddress())
	k.SetPool(ctx, pool)

	share := sdk.NewCoin(pool.GetPoolShareDenom(), initialTotalShare)
	err = k.mintLpShares(ctx, share, fromAddress)
	if err != nil {
		return err
	}
	return nil
}
