package types_test

import (
	"github.com/cosmos/cosmos-sdk/x/amm/types"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	poolToken0         = "pool/0"
	poolToken1         = "pool/1"
	poolToken2         = "pool/2"
	amount1            = sdk.NewInt(100)
	amount2            = sdk.NewInt(20)
	validCoinShares    = sdk.NewCoins(sdk.NewCoin(poolToken1, sdk.OneInt()))
	unsortedCoinShares = []sdk.Coin{sdk.NewCoin(poolToken1, sdk.OneInt()), sdk.NewCoin(poolToken0, sdk.OneInt())}
	zeroCoinShares     = []sdk.Coin{sdk.NewCoin(poolToken1, sdk.ZeroInt())}
)

func TestPoolShareAdd(t *testing.T) {
	tests := []struct {
		title         string
		share         sdk.Int
		toAdd         sdk.Coin
		expectedValue sdk.Coins
	}{
		{"adding existing shares", amount1, sdk.NewCoin(poolToken1, amount2), sdk.NewCoins(sdk.NewCoin(poolToken1, sdk.NewInt(120)))},
		{"adding different shares", amount1, sdk.NewCoin(poolToken2, amount2), sdk.NewCoins(sdk.NewCoin(poolToken1, amount1), sdk.NewCoin(poolToken2, amount2))},
		{"adding different shares with sorting", amount1, sdk.NewCoin(poolToken0, amount2), sdk.Coins([]sdk.Coin{{poolToken0, amount2}, {poolToken1, amount1}})},
	}
	for i, tc := range tests {
		poolShare := types.NewPoolShare(account1, sdk.NewCoin(poolToken1, tc.share))
		poolShare.AddShare(tc.toAdd)
		require.Equal(t, tc.expectedValue, poolShare.Shares, "test: (%v) %s", i, tc.title)
	}
}

func TestPoolShareSubtract(t *testing.T) {
	tests := []struct {
		title         string
		shares        sdk.Coins
		toSubtract    sdk.Coin
		expectedValue sdk.Coins
	}{
		{"subtract", sdk.NewCoins(sdk.NewCoin(poolToken1, amount1)), sdk.NewCoin(poolToken1, amount2), sdk.NewCoins(sdk.NewCoin(poolToken1, amount1.Sub(amount2)))},
		{"remove 0 values", sdk.NewCoins(sdk.NewCoin(poolToken1, amount1), sdk.NewCoin(poolToken2, amount2)), sdk.NewCoin(poolToken2, amount2), sdk.NewCoins(sdk.NewCoin(poolToken1, amount1))},
	}
	for i, tc := range tests {
		poolShare := types.PoolShare{account1.String(), tc.shares}
		poolShare.SubtractShare(tc.toSubtract)
		require.Equal(t, tc.expectedValue, poolShare.Shares, "test: (%v) %s", i, tc.title)
	}
}

func TestPoolShare(t *testing.T) {
	tests := []struct {
		title      string
		address    sdk.AccAddress
		shares     sdk.Coins
		expectPass bool
	}{
		{"nil From address", nil, validCoinShares, false},
		{"unsorted shares", account1, unsortedCoinShares, false},
		{"zero shares", account1, zeroCoinShares, false},
		{"valid test case", account1, validCoinShares, true},
	}
	for i, tc := range tests {
		poolShare := types.PoolShare{tc.address.String(), tc.shares}
		if tc.expectPass {
			require.NoError(t, poolShare.Validate(), "test: (%v) %s", i, tc.title)
		} else {
			require.Error(t, poolShare.Validate(), "test: (%v) %s", i, tc.title)
		}
	}
}
