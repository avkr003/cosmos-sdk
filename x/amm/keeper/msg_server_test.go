package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/cosmos/cosmos-sdk/x/amm/types"
	"testing"
)

func (s *KeeperTestSuite) TestMsgUpdateParams() {
	ctx, keeper := s.ctx, s.ammKeeper
	require := s.Require()

	testCases := []struct {
		name      string
		input     *ammtypes.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "valid params",
			input: &ammtypes.MsgUpdateParams{
				Authority: keeper.GetAuthority(),
				Params:    ammtypes.DefaultParams(),
			},
			expErr: false,
		},
		{
			name: "invalid authority",
			input: &ammtypes.MsgUpdateParams{
				Authority: "invalid",
				Params:    ammtypes.DefaultParams(),
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "negative swap fee",
			input: &ammtypes.MsgUpdateParams{
				Authority: keeper.GetAuthority(),
				Params: ammtypes.Params{
					MaxSwapFee:        math.LegacyNewDec(-10),
					SwapAllowedTokens: []string{"ubtc"},
				},
			},
			expErr:    true,
			expErrMsg: "max fee should be between 0 and 1",
		},
		{
			name: "zero swap fee",
			input: &ammtypes.MsgUpdateParams{
				Authority: keeper.GetAuthority(),
				Params: ammtypes.Params{
					MaxSwapFee:        sdk.ZeroDec(),
					SwapAllowedTokens: []string{"ubtc"},
				},
			},
			expErr:    true,
			expErrMsg: "max fee should be between 0 and 1",
		},
		{
			name: "one swap fee",
			input: &ammtypes.MsgUpdateParams{
				Authority: keeper.GetAuthority(),
				Params: ammtypes.Params{
					MaxSwapFee:        sdk.OneDec(),
					SwapAllowedTokens: []string{"ubtc"},
				},
			},
			expErr:    true,
			expErrMsg: "max fee should be between 0 and 1",
		},
		{
			name: "more than 1 swap fee",
			input: &ammtypes.MsgUpdateParams{
				Authority: keeper.GetAuthority(),
				Params: ammtypes.Params{
					MaxSwapFee:        math.LegacyNewDec(20),
					SwapAllowedTokens: []string{"ubtc"},
				},
			},
			expErr:    true,
			expErrMsg: "max fee should be between 0 and 1",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := keeper.UpdateParams(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
				require.Contains(err.Error(), tc.expErrMsg)
			} else {
				require.NoError(err)
			}
		})
	}
}
