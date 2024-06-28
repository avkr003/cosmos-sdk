package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePool = "createPool"
	TypeMsgJoinPool   = "joinPool"
	TypeMsgSwap       = "swap"
	TypeMsgExitPool   = "exitPool"
)

var (
	_ sdk.Msg = &MsgCreatePoolMessage{}
	_ sdk.Msg = &MsgJoinPoolMessage{}
	_ sdk.Msg = &MsgSwapMessage{}
	_ sdk.Msg = &MsgExitPoolMessage{}
)

func NewMsgCreatePoolMessage(from sdk.AccAddress, token1, token2 sdk.Coin, fee sdk.Dec) MsgCreatePoolMessage {

	return MsgCreatePoolMessage{
		From:   from.String(),
		Token1: token1,
		Token2: token2,
		Fee:    fee,
	}
}

func (msg MsgCreatePoolMessage) Route() string { return RouterKey }

func (msg MsgCreatePoolMessage) Type() string { return TypeMsgCreatePool }

func (msg MsgCreatePoolMessage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.From)}
}

func (msg MsgCreatePoolMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreatePoolMessage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	err = msg.Token1.Validate()
	if err != nil {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid coin: %s", err)
	}

	err = msg.Token2.Validate()
	if err != nil {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid coin: %s", err)
	}

	if msg.Fee.LTE(sdk.ZeroDec()) {
		return ErrInvalidFees
	}
	return nil
}

func NewMsgJoinPoolMessage(from sdk.AccAddress, poolId uint64, token sdk.Coin) MsgJoinPoolMessage {

	return MsgJoinPoolMessage{
		From:   from.String(),
		PoolId: poolId,
		Token:  token,
	}
}

func (msg MsgJoinPoolMessage) Route() string { return RouterKey }

func (msg MsgJoinPoolMessage) Type() string { return TypeMsgJoinPool }

func (msg MsgJoinPoolMessage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.From)}
}

func (msg MsgJoinPoolMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgJoinPoolMessage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	err = msg.Token.Validate()
	if err != nil {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid coin: %s", err)
	}

	return nil
}

func NewMsgSwapMessage(from sdk.AccAddress, poolId uint64, token sdk.Coin) MsgSwapMessage {

	return MsgSwapMessage{
		From:   from.String(),
		PoolId: poolId,
		Token:  token,
	}
}

func (msg MsgSwapMessage) Route() string { return RouterKey }

func (msg MsgSwapMessage) Type() string { return TypeMsgSwap }

func (msg MsgSwapMessage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.From)}
}

func (msg MsgSwapMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSwapMessage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	err = msg.Token.Validate()
	if err != nil {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid coin: %s", err)
	}

	return nil
}

func NewMsgExitPoolMessage(from sdk.AccAddress, poolId uint64, lpShare sdk.Dec, withdrawAll bool) MsgExitPoolMessage {

	return MsgExitPoolMessage{
		From:        from.String(),
		PoolId:      poolId,
		LpShare:     lpShare,
		WithdrawAll: withdrawAll,
	}
}

func (msg MsgExitPoolMessage) Route() string { return RouterKey }

func (msg MsgExitPoolMessage) Type() string { return TypeMsgExitPool }

func (msg MsgExitPoolMessage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.From)}
}

func (msg MsgExitPoolMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgExitPoolMessage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if msg.LpShare.LTE(sdk.ZeroDec()) {
		return ErrInvalidLpShares
	}
	return nil
}
