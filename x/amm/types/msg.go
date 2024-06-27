package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePool = "createPool"
)

var (
	_ sdk.Msg = &MsgCreatePoolMessage{}
)

func NewMsgCreatePoolMessage(from sdk.AccAddress, token1, token2 sdk.Coin, fee sdk.Dec) MsgCreatePoolMessage {

	return MsgCreatePoolMessage{
		FromAddress: from.String(),
		Token_1:     token1,
		Token_2:     token2,
		Fee:         fee,
	}
}

func (msg MsgCreatePoolMessage) Route() string { return RouterKey }

func (msg MsgCreatePoolMessage) Type() string { return TypeMsgCreatePool }

func (msg MsgCreatePoolMessage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.FromAddress)}
}

func (msg MsgCreatePoolMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreatePoolMessage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	err = msg.Token_1.Validate()
	if err != nil {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid coin: %s", err)
	}

	err = msg.Token_2.Validate()
	if err != nil {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid coin: %s", err)
	}

	if msg.Fee.LTE(sdk.ZeroDec()) {
		return ErrInvalidFees
	}
	return nil
}
