package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *PoolShare) GetAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(p.Address)
}

func (p *PoolShare) GetShare(denom string) (sdk.DecCoin, int) {
	for i, decCoin := range p.GetShares() {
		if decCoin.Denom == denom {
			return decCoin, i
		}
	}
	return sdk.DecCoin{}, -1
}

func (p *PoolShare) AddShare(toAdd sdk.DecCoin) {
	coin, found := p.GetShare(toAdd.Denom)
	if found == -1 {
		p.Shares = append(p.Shares, toAdd)
		p.Shares = p.Shares.Sort()
	} else {
		p.Shares[found] = coin.Add(toAdd)
	}
}

func (p *PoolShare) SubtractShare(toSubtract sdk.DecCoin) error {
	coin, found := p.GetShare(toSubtract.Denom)
	if found != -1 {
		value := coin.Sub(toSubtract)
		if value.Amount.LT(sdk.ZeroDec()) {
			return ErrPoolShareGreater
		}
		if value.Equal(sdk.ZeroDec()) {
			p.Shares = append(p.Shares[:found], p.Shares[found+1:]...)
		} else {
			p.Shares[found] = value
		}
		return p.Shares[found].Validate()
	}
	return ErrLpSharesNotFound
}

func (p PoolShare) GetKey() []byte {
	return GetPoolShareKey(p.GetAccAddress())
}

func NewPoolShare(account sdk.AccAddress, share sdk.DecCoin) PoolShare {
	return PoolShare{
		Address: account.String(),
		Shares:  sdk.DecCoins{share},
	}
}

func MustMarshalPoolShare(cdc codec.BinaryCodec, poolShare *PoolShare) []byte {
	return cdc.MustMarshal(poolShare)
}

func MustUnmarshalPoolShare(cdc codec.BinaryCodec, value []byte) (poolShare PoolShare) {
	err := cdc.Unmarshal(value, &poolShare)
	if err != nil {
		panic(err)
	}

	return poolShare
}
