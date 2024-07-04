package types_test

import (
	"github.com/cosmos/cosmos-sdk/x/amm/types"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	validShares    = sdk.OneInt()
	negativeShares = sdk.OneInt().MulRaw(-1)
)

func TestPool(t *testing.T) {
	tests := []struct {
		title          string
		creatorAddress sdk.AccAddress
		token1         sdk.Coin
		token2         sdk.Coin
		fee            sdk.Dec
		expectPass     bool
		poolName       string
		override       bool
	}{
		{"nil From address", nil, coin1, coin2, validFee, false, "", false},
		{"negative coin", account1, coin1, negativeCoin, validFee, false, "", false},
		{"invalid denom coin", account1, invalidDenomCoin, coin2, validFee, false, "", false},
		{"same denom coin", account1, coin2, coin2, validFee, false, "", false},
		{"zero amount coins", account1, zeroCoin, coin2, validFee, false, "", false},
		{"fee greater than or equal to 1 ", account1, coin1, coin2, oneInvalidFee, false, "", false},
		{"fee equal to 0 ", account1, coin1, coin2, zeroInvalidFee, false, "", false},
		{"fee less than  0 ", account1, coin1, coin2, negativeInvalidFee, false, "", false},
		{"total shares less than  0 ", account1, coin1, coin2, validFee, false, "", false},
		{"total shares 0 with non zero token ", account1, coin1, coin2, validFee, false, "", false},
		{"empty pool name", account1, coin1, coin2, validFee, false, "", true},
		{"invalid pool name", account1, coin1, coin2, validFee, false, "pool/1/uabc/ubcd", true},
		{"valid test case", account1, coin1, coin2, validFee, true, "", false},
	}
	for i, tc := range tests {
		pool := types.NewPool(1, tc.token1, tc.token2, tc.fee, tc.creatorAddress)
		if tc.override {
			pool = types.Pool{
				Id:      1,
				Name:    tc.poolName,
				Token1:  tc.token1,
				Token2:  tc.token2,
				Fee:     tc.fee,
				Creator: tc.creatorAddress.String(),
			}
		}
		if tc.expectPass {
			require.NoError(t, pool.Validate(), "test: (%v) %s", i, tc.title)
		} else {
			require.Error(t, pool.Validate(), "test: (%v) %s", i, tc.title)
		}
	}
}
