package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sigs.k8s.io/yaml"
)

var DefaultSwapAllowedTokens = []string{"uelys", "ubtc", "ueth", "uatom", "uusd"}

var DefaultMaxSwapFeee = sdk.MustNewDecFromStr("0.5")

func NewParams(allowedTokens []string, maxFee sdk.Dec) Params {
	return Params{
		SwapAllowedTokens: allowedTokens,
		MaxSwapFee:        maxFee,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultSwapAllowedTokens, DefaultMaxSwapFeee)
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func MustUnmarshalParams(cdc *codec.LegacyAmino, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}

	return params
}

func UnmarshalParams(cdc *codec.LegacyAmino, value []byte) (params Params, err error) {
	err = cdc.Unmarshal(value, &params)
	if err != nil {
		return
	}

	return
}

func (p Params) Validate() error {
	err := validateAllowedDenoms(p.SwapAllowedTokens)
	if err != nil {
		return err
	}

	err = validateMaxFee(p.MaxSwapFee)
	if err != nil {
		return err
	}

	return nil
}

func validateAllowedDenoms(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, denom := range v {
		if err := sdk.ValidateDenom(denom); err != nil {
			return err
		}
	}
	return nil
}

func validateMaxFee(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.LTE(sdk.ZeroDec()) || v.GTE(sdk.OneDec()) {
		return ErrInvalidFees.Wrapf("max fee should be between 0 and 1")
	}
	return nil
}
