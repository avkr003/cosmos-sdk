package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"sigs.k8s.io/yaml"
)

const DefaultFeeDistributionPeriod = 10000

func NewParams(feeDistributionPeriod uint32) Params {
	return Params{
		FeeDistributionPeriod: feeDistributionPeriod,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultFeeDistributionPeriod)
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
	return nil
}
