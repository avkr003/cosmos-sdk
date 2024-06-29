package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sigs.k8s.io/yaml"
)

var DefaultAllowedTokens = []string{"uelys", "ubtc", "ueth", "uatom", "uusd"}

func NewParams(allowedTokens []string) Params {
	return Params{
		AllowedTokens: allowedTokens,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultAllowedTokens)
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
	for _, denom := range p.AllowedTokens {
		if err := sdk.ValidateDenom(denom); err != nil {
			return err
		}
	}
	return nil
}
