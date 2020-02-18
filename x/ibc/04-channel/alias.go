package channel

// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/cosmos/cosmos-sdk/x/ibc/04-channel/keeper
// ALIASGEN: github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types

import (
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/keeper"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
)

const (
	SubModuleName            = types.SubModuleName
	StoreKey                 = types.StoreKey
	RouterKey                = types.RouterKey
	QuerierRoute             = types.QuerierRoute
	QueryAllChannels         = types.QueryAllChannels
	QueryChannel             = types.QueryChannel
)

var (
	// functions aliases
	NewKeeper                    = keeper.NewKeeper
	QuerierChannels              = keeper.QuerierChannels
	NewChannel                   = types.NewChannel
	NewCounterparty              = types.NewCounterparty
	RegisterCodec                = types.RegisterCodec
	ErrChannelExists             = types.ErrChannelExists
	ErrChannelNotFound           = types.ErrChannelNotFound
	ErrInvalidCounterparty       = types.ErrInvalidCounterparty
	ErrChannelCapabilityNotFound = types.ErrChannelCapabilityNotFound
	ErrInvalidPacket             = types.ErrInvalidPacket
	ErrSequenceSendNotFound      = types.ErrSequenceSendNotFound
	ErrSequenceReceiveNotFound   = types.ErrSequenceReceiveNotFound
	ErrPacketTimeout             = types.ErrPacketTimeout
	ErrInvalidChannel            = types.ErrInvalidChannel
	ErrInvalidChannelState       = types.ErrInvalidChannelState
	NewMsgChannelOpenInit        = types.NewMsgChannelOpenInit
	NewMsgChannelOpenTry         = types.NewMsgChannelOpenTry
	NewMsgChannelOpenAck         = types.NewMsgChannelOpenAck
	NewMsgChannelOpenConfirm     = types.NewMsgChannelOpenConfirm
	NewMsgChannelCloseInit       = types.NewMsgChannelCloseInit
	NewMsgChannelCloseConfirm    = types.NewMsgChannelCloseConfirm
	NewMsgPacket                 = types.NewMsgPacket
	NewMsgTimeout                = types.NewMsgTimeout
	NewMsgAcknowledgement        = types.NewMsgAcknowledgement
	NewPacket                    = types.NewPacket
	NewChannelResponse           = types.NewChannelResponse

	// variable aliases
	SubModuleCdc                 = types.SubModuleCdc
	EventTypeChannelOpenInit     = types.EventTypeChannelOpenInit
	EventTypeChannelOpenTry      = types.EventTypeChannelOpenTry
	EventTypeChannelOpenAck      = types.EventTypeChannelOpenAck
	EventTypeChannelOpenConfirm  = types.EventTypeChannelOpenConfirm
	EventTypeChannelCloseInit    = types.EventTypeChannelCloseInit
	EventTypeChannelCloseConfirm = types.EventTypeChannelCloseConfirm
	AttributeValueCategory       = types.AttributeValueCategory
)

type (
	Keeper                 = keeper.Keeper
	Channel                = types.Channel
	Counterparty           = types.Counterparty
	ClientKeeper           = types.ClientKeeper
	ConnectionKeeper       = types.ConnectionKeeper
	PortKeeper             = types.PortKeeper
	MsgChannelOpenInit     = types.MsgChannelOpenInit
	MsgChannelOpenTry      = types.MsgChannelOpenTry
	MsgChannelOpenAck      = types.MsgChannelOpenAck
	MsgChannelOpenConfirm  = types.MsgChannelOpenConfirm
	MsgChannelCloseInit    = types.MsgChannelCloseInit
	MsgChannelCloseConfirm = types.MsgChannelCloseConfirm
	MsgPacket              = types.MsgPacket
	MsgAcknowledgement     = types.MsgAcknowledgement
	MsgTimeout             = types.MsgTimeout
	Packet                 = types.Packet
	ChannelResponse        = types.ChannelResponse
)
