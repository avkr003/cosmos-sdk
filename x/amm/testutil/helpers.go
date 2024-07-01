package testutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"testing"
)

// Helper is a structure which wraps the staking message server
// and provides methods useful in tests
type Helper struct {
	t       *testing.T
	msgSrvr stakingtypes.MsgServer
	k       *keeper.Keeper

	Ctx        sdk.Context
	Commission stakingtypes.CommissionRates
	// Coin Denomination
	Denom string
}

// NewHelper creates a new instance of Helper.
//func NewHelper(t *testing.T, ctx sdk.Context, k *keeper.Keeper) *Helper {
//	return &Helper{t, keeper.NewMsgServerImpl(k), k, ctx, ZeroCommission(), sdk.DefaultBondDenom}
//}
