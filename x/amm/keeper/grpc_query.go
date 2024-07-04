package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/amm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Parameters(c context.Context, req *types.QueryParametersRequest) (*types.QueryParametersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &types.QueryParametersResponse{Params: params}, nil
}

func (k Keeper) Pool(c context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "pool not found")
	}
	return &types.QueryPoolResponse{Pool: pool}, nil
}

func (k Keeper) PoolShares(c context.Context, req *types.QueryPoolSharesRequest) (*types.QueryPoolSharesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address: %s", err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	allBalances := k.bankKeeper.GetAllBalances(ctx, address)
	pools := k.GetAllPools(ctx)

	balances := sdk.NewCoins()
	for _, pool := range pools {
		amount := allBalances.AmountOf(pool.GetPoolShareDenom())
		if !amount.IsZero() {
			balances = balances.Add(sdk.NewCoin(pool.GetPoolShareDenom(), amount))
		}
	}

	return &types.QueryPoolSharesResponse{Shares: balances}, nil
}

func (k Keeper) Pools(c context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	pools := []types.Pool{}
	ctx := sdk.UnwrapSDKContext(c)

	poolStore := k.getPoolStore(ctx)
	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key, value []byte) error {
		pool := types.MustUnmarshalPool(k.cdc, value)
		pools = append(pools, pool)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryPoolsResponse{Pools: pools, Pagination: pageRes}, nil
}
