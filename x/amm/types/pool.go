package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

func (p Pool) GetFee() sdk.Dec {
	return p.Fee
}

func (p *Pool) SetFee(v sdk.Dec) {
	p.Fee = v
}

func (p Pool) GetTotalShares() sdk.Dec {
	return p.TotalShares
}

func (p *Pool) setTotalShares(v sdk.Dec) {
	p.TotalShares = v
}

func (p *Pool) AddTotalShares(v sdk.Dec) {
	p.TotalShares = p.TotalShares.Add(v)
}

func (p Pool) GetCreatorAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(p.Creator)
}

func (p Pool) GetPoolAddress() sdk.AccAddress {
	key := append([]byte("pool"), sdk.Uint64ToBigEndian(p.Id)...)
	return address.Module(ModuleName, key)
}

func (p Pool) GetPoolShareDenom() string {
	return fmt.Sprintf("pool/%d", p.GetId())
}

func (p Pool) GetKey() []byte {
	return GetPoolKey(p.GetId())
}

func NewPool(id uint64, token1, token2 sdk.Coin, fee sdk.Dec, creator sdk.AccAddress, totalShare sdk.Dec) Pool {
	return Pool{
		Id:          id,
		Token1:      token1,
		Token2:      token2,
		Fee:         fee,
		Creator:     creator.String(),
		TotalShares: totalShare,
	}
}

func MustMarshalPool(cdc codec.BinaryCodec, pool *Pool) []byte {
	return cdc.MustMarshal(pool)
}

func MustUnmarshalPool(cdc codec.BinaryCodec, value []byte) (pool Pool) {
	err := cdc.Unmarshal(value, &pool)
	if err != nil {
		panic(err)
	}

	return pool
}
