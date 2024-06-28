package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/amm/types"
	"strconv"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) CreatePool(goCtx context.Context, msg *types.MsgCreatePoolMessage) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	creatorAddress := sdk.MustAccAddressFromBech32(msg.From)

	pool, err := k.createNewPool(ctx, creatorAddress, msg.Token1, msg.Token2, msg.Fee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventNewPool,
			sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(pool.GetId(), 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.From),
			sdk.NewAttribute(types.AttributeKeyTotalShares, pool.TotalShares.String()),
		),
	})

	return &types.MsgCreatePoolResponse{}, nil
}

func (k Keeper) JoinPool(goCtx context.Context, msg *types.MsgJoinPoolMessage) (*types.MsgJoinPoolResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	//fromAddress := sdk.MustAccAddressFromBech32(msg.From)
	//

	return &types.MsgJoinPoolResponse{}, nil
}

func (k Keeper) Swap(goCtx context.Context, msg *types.MsgSwapMessage) (*types.MsgSwapResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	//fromAddress := sdk.MustAccAddressFromBech32(msg.From)

	return &types.MsgSwapResponse{}, nil
}

func (k Keeper) ExitPool(goCtx context.Context, msg *types.MsgExitPoolMessage) (*types.MsgExitPoolResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	//fromAddress := sdk.MustAccAddressFromBech32(msg.From)

	return &types.MsgExitPoolResponse{}, nil
}
