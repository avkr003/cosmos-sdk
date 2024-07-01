package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *PoolShare) GetAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(p.Address)
}

func (p *PoolShare) GetShare(denom string) sdk.Int {
	return p.GetShares().AmountOf(denom)
}

func (p *PoolShare) AddShare(toAdd sdk.Coin) {
	p.Shares = p.Shares.Add(toAdd)
}

func (p *PoolShare) SubtractShare(toSubtract sdk.Coin) {
	p.Shares = p.Shares.Sub(toSubtract)
}

func (p PoolShare) GetKey() []byte {
	return GetPoolShareKey(p.GetAccAddress())
}

func (p PoolShare) Validate() error {
	if _, err := sdk.AccAddressFromBech32(p.Address); err != nil {
		return err
	}
	if err := p.Shares.Validate(); err != nil {
		return err
	}
	for _, share := range p.GetShares() {
		if share.Amount.LTE(sdk.ZeroInt()) {
			return ErrInvalidLpShares
		}
	}
	return nil
}

func NewPoolShare(account sdk.AccAddress, share sdk.Coin) PoolShare {
	return PoolShare{
		Address: account.String(),
		Shares:  sdk.NewCoins(share),
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
