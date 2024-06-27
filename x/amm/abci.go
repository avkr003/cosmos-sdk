package amm

import (
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/amm/keeper"
)

func BeginBlocker(ctx sdk.Context, block abci.RequestBeginBlock, keeper keeper.Keeper) {
	fmt.Println("Implement Begin Blocker")
}
