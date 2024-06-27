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
	creatorAddress := sdk.MustAccAddressFromBech32(msg.FromAddress)
	initialShare := sdk.NewInt(1000000)
	poolId := k.GetNextPoolNumber(ctx)

	pool := types.NewPool(poolId, msg.Token_1, msg.Token_2, msg.Fee, creatorAddress, initialShare)
	poolShare := types.NewPoolShare(poolId, creatorAddress, initialShare)

	err := k.createNewPool(ctx, pool, poolShare)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventNewPool,
			sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.FromAddress),
		),
	})

	return &types.MsgCreatePoolResponse{}, nil
}
