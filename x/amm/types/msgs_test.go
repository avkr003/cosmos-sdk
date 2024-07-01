package types_test

import (
	"github.com/cosmos/cosmos-sdk/x/amm/types"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	coin1                 = sdk.NewInt64Coin("ubtc", 1000)
	coin2                 = sdk.NewInt64Coin("uusd", 1000)
	zeroCoin              = sdk.NewInt64Coin("ueth", 0)
	coin4                 = sdk.NewInt64Coin("uatom", 1000)
	negativeCoin          = sdk.Coin{"ueth", sdk.NewInt(-100)}
	invalidDenomCoin      = sdk.Coin{"u", sdk.NewInt(100)}
	validCoins            = sdk.NewCoins(coin1, coin2)
	unsortedInvalidCoins  = sdk.Coins{coin2, coin1}
	zeroInvalidCoins      = sdk.Coins{coin1, zeroCoin}
	lengthInvalidCoins    = sdk.Coins{coin1, coin4, coin2}
	account1              = sdk.MustAccAddressFromBech32("cosmos12wrytpuqzrp2gp0qtcad9q2033ft7q0dryqycg")
	validFee, _           = sdk.NewDecFromStr("0.03")
	zeroInvalidFee        = sdk.ZeroDec()
	oneInvalidFee         = sdk.OneDec()
	negativeInvalidFee, _ = sdk.NewDecFromStr("-0.03")
)

func TestMsgCreatePoolMessage(t *testing.T) {
	tests := []struct {
		title       string
		fromAddress sdk.AccAddress
		tokens      sdk.Coins
		fee         sdk.Dec
		expectPass  bool
	}{
		{"nil From address", nil, validCoins, validFee, false},
		{"more than 2 coins", account1, lengthInvalidCoins, validFee, false},
		{"unsorted coins", account1, unsortedInvalidCoins, validFee, false},
		{"zero amount coins", account1, zeroInvalidCoins, validFee, false},
		{"fee greater than or equal to 1 ", account1, validCoins, oneInvalidFee, false},
		{"fee equal to 0 ", account1, validCoins, zeroInvalidFee, false},
		{"fee less than  0 ", account1, validCoins, negativeInvalidFee, false},
		{"valid test case", account1, validCoins, validFee, true},
	}
	for i, tc := range tests {
		msg := types.NewMsgCreatePoolMessage(tc.fromAddress, tc.tokens, tc.fee)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: (%v) %s", i, tc.title)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: (%v) %s", i, tc.title)
		}
	}
}

func TestMsgJoinPoolMessage(t *testing.T) {
	tests := []struct {
		title       string
		fromAddress sdk.AccAddress
		poolId      uint64
		tokens      sdk.Coins
		expectPass  bool
	}{
		{"nil From address", nil, 1, validCoins, false},
		{"more than 2 coins", account1, 1, lengthInvalidCoins, false},
		{"unsorted coins", account1, 1, unsortedInvalidCoins, false},
		{"zero amount coins", account1, 1, zeroInvalidCoins, false},
		{"valid test case", account1, 1, validCoins, true},
	}
	for i, tc := range tests {
		msg := types.NewMsgJoinPoolMessage(tc.fromAddress, tc.poolId, tc.tokens)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: (%v) %s", i, tc.title)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: (%v) %s", i, tc.title)
		}
	}
}

func TestMsgSwapMessage(t *testing.T) {
	tests := []struct {
		title       string
		fromAddress sdk.AccAddress
		poolId      uint64
		token       sdk.Coin
		expectPass  bool
	}{
		{"nil From address", nil, 1, coin1, false},
		{"negative coin", account1, 1, negativeCoin, false},
		{"zero amount coin", account1, 1, zeroCoin, false},
		{"valid test case", account1, 1, coin1, true},
	}
	for i, tc := range tests {
		msg := types.NewMsgSwapMessage(tc.fromAddress, tc.poolId, tc.token)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: (%v) %s", i, tc.title)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: (%v) %s", i, tc.title)
		}
	}
}

func TestMsgExitPoolMessage(t *testing.T) {
	tests := []struct {
		title       string
		fromAddress sdk.AccAddress
		poolId      uint64
		lpShare     sdk.Int
		expectPass  bool
	}{
		{"nil From address", nil, 1, sdk.OneInt(), false},
		{"negative shares", account1, 1, sdk.OneInt().MulRaw(-1), false},
		{"zero shares", account1, 1, sdk.ZeroInt(), false},
		{"valid test case", account1, 1, sdk.OneInt(), true},
	}
	for i, tc := range tests {
		msg := types.NewMsgExitPoolMessage(tc.fromAddress, tc.poolId, tc.lpShare, false)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: (%v) %s", i, tc.title)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: (%v) %s", i, tc.title)
		}
	}
}
