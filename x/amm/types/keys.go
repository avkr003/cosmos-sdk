package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	ModuleName = "amm"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	PoolKey       = []byte{0x11}
	PoolShareKey  = []byte{0x12}
	PoolNumberKey = []byte{0x13}
	ParamsKey     = []byte{0x14}
)

func GetPoolKey(id uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return append(PoolKey, b...)
}

func GetPoolShareKey(accountAddress sdk.AccAddress) []byte {
	return append(PoolShareKey, address.MustLengthPrefix(accountAddress)...)
}

func GetPoolSharesKey(id uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return append(PoolShareKey, b...)
}
