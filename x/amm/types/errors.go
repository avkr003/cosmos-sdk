package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidFees = sdkerrors.Register(ModuleName, 1, "invalid pool fees")
)
