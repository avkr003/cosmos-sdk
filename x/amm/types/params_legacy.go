package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyMaxFee        = []byte("MaxSwapFee")
	KeyAllowedTokens = []byte("SwapAllowedTokens")
)

var _ paramtypes.ParamSet = &Params{}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxFee, &p.MaxSwapFee, validateMaxFee),
		paramtypes.NewParamSetPair(KeyAllowedTokens, &p.SwapAllowedTokens, validateAllowedDenoms),
	}
}
