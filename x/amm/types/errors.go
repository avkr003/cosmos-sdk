package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidFees      = sdkerrors.Register(ModuleName, 1, "invalid pool fees")
	ErrInvalidLpShares  = sdkerrors.Register(ModuleName, 2, "invalid lp shares fees")
	ErrInvalidPoolId    = sdkerrors.Register(ModuleName, 3, "invalid poolId")
	ErrLpSharesNotFound = sdkerrors.Register(ModuleName, 4, "Lp shares not found")
)
