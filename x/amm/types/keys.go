package types

import (
	"encoding/binary"
)

const (
	ModuleName = "amm"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	PoolKey       = []byte{0x11}
	PoolNumberKey = []byte{0x12}
	ParamsKey     = []byte{0x13}
)

func GetPoolKey(id uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return append(PoolKey, b...)
}
