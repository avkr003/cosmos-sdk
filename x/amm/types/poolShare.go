package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *PoolShare) GetAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(p.Address)
}

func (p *PoolShare) GetShares() sdk.Int {
	return p.Shares
}

func (p *PoolShare) AddShares(v sdk.Int) {
	p.Shares = p.Shares.Add(v)
}

func (p PoolShare) GetKey() []byte {
	return GetPoolShareKey(p.GetId(), p.GetAccAddress())
}

func NewPoolShare(id uint64, account sdk.AccAddress, shares sdk.Int) PoolShare {
	return PoolShare{
		Id:      id,
		Address: account.String(),
		Shares:  shares,
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
