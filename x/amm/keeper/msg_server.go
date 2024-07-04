package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/amm/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"strconv"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) CreatePool(goCtx context.Context, msg *types.MsgCreatePoolMessage) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	creatorAddress := sdk.MustAccAddressFromBech32(msg.From)

	pool, err := k.createNewPool(ctx, creatorAddress, msg.Tokens, msg.Fee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventNewPool,
			sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(pool.GetId(), 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.From),
			sdk.NewAttribute(types.AttributeKeyPoolAddress, pool.GetPoolAddress().String()),
		),
	})

	return &types.MsgCreatePoolResponse{}, nil
}

func (k Keeper) JoinPool(goCtx context.Context, msg *types.MsgJoinPoolMessage) (*types.MsgJoinPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fromAddress := sdk.MustAccAddressFromBech32(msg.From)

	pool, err := k.joinPool(ctx, msg.PoolId, fromAddress, msg.Tokens)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventLiquidityAdded,
			sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(pool.GetId(), 10)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.From),
		),
	})

	return &types.MsgJoinPoolResponse{}, nil
}

func (k Keeper) Swap(goCtx context.Context, msg *types.MsgSwapMessage) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fromAddress := sdk.MustAccAddressFromBech32(msg.From)

	tokenOut, err := k.swap(ctx, msg.PoolId, fromAddress, msg.Token)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventSwap,
			sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.From),
			sdk.NewAttribute(types.AttributeKeyTokenIn, msg.Token.String()),
			sdk.NewAttribute(types.AttributeKeyTokenOut, tokenOut.String()),
		),
	})

	return &types.MsgSwapResponse{}, nil
}

func (k Keeper) ExitPool(goCtx context.Context, msg *types.MsgExitPoolMessage) (*types.MsgExitPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fromAddress := sdk.MustAccAddressFromBech32(msg.From)

	err := k.exitPool(ctx, msg.PoolId, fromAddress, msg.LpShare, msg.WithdrawAll)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventLiquidityRemoved,
			sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.From),
		),
	})

	return &types.MsgExitPoolResponse{}, nil
}

func (k Keeper) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
