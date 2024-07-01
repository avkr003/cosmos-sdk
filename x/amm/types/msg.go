package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePool   = "createPool"
	TypeMsgJoinPool     = "joinPool"
	TypeMsgSwap         = "swap"
	TypeMsgExitPool     = "exitPool"
	TypeMsgUpdateParams = "update_params"
)

var (
	_ sdk.Msg = &MsgCreatePoolMessage{}
	_ sdk.Msg = &MsgJoinPoolMessage{}
	_ sdk.Msg = &MsgSwapMessage{}
	_ sdk.Msg = &MsgExitPoolMessage{}
	_ sdk.Msg = &MsgUpdateParams{}
)

func NewMsgCreatePoolMessage(from sdk.AccAddress, tokens sdk.Coins, fee sdk.Dec) MsgCreatePoolMessage {

	return MsgCreatePoolMessage{
		From:   from.String(),
		Tokens: tokens,
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

	if len(msg.Tokens) != 2 {
		return ErrInvalidTokens.Wrapf("only 2 tokens has to be given")
	}

	err = msg.Tokens.Validate()
	if err != nil {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid coin: %s", err)
	}

	if msg.Fee.LTE(sdk.ZeroDec()) {
		return ErrInvalidFees.Wrapf("fees less than 0")
	}

	if msg.Fee.GTE(sdk.OneDec()) {
		return ErrInvalidFees.Wrapf("fees greater than 1")
	}
	return nil
}

func NewMsgJoinPoolMessage(from sdk.AccAddress, poolId uint64, tokens sdk.Coins) MsgJoinPoolMessage {

	return MsgJoinPoolMessage{
		From:   from.String(),
		PoolId: poolId,
		Tokens: tokens,
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

	if len(msg.Tokens) == 0 || len(msg.Tokens) > 2 {
		return ErrInvalidTokens.Wrapf("only 1 or 2 token can be given")
	}

	err = msg.Tokens.Validate()
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

	if msg.Token.Amount.IsZero() {
		return sdkerrors.ErrInvalidCoins.Wrapf("cannot swap 0 tokens: %s", msg.Token.String())
	}

	return nil
}

func NewMsgExitPoolMessage(from sdk.AccAddress, poolId uint64, lpShare sdk.Int, withdrawAll bool) MsgExitPoolMessage {

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

	if msg.LpShare.LTE(sdk.ZeroInt()) {
		return ErrInvalidLpShares
	}
	return nil
}

func (m *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}
	return m.Params.Validate()
}

func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}
