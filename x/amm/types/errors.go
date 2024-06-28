package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidFees              = sdkerrors.Register(ModuleName, 1, "invalid pool fees")
	ErrInvalidLpShares          = sdkerrors.Register(ModuleName, 2, "invalid lp shares fees")
	ErrLpSharesNotFound         = sdkerrors.Register(ModuleName, 4, "Lp shares not found")
	ErrPoolNotFound             = sdkerrors.Register(ModuleName, 5, "Pool not found")
	ErrInvalidToken             = sdkerrors.Register(ModuleName, 6, "Invalid token for the pool")
	ErrPoolShareGreater         = sdkerrors.Register(ModuleName, 7, "Pool share value greater than total")
	ErrRedeemingMoreThanAllowed = sdkerrors.Register(ModuleName, 8, "Redeeming pool share than allowed")
)
