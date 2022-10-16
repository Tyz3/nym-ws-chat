package chat_payload

import (
	. "nym-ws-chat/client/web_socket_packet"
)

type payload struct {
	Sig    byte
	Length uint64
}

type PayloadReader interface {
	ReadFrom(reader *WSPacketReader)
	String() string
}

type PayloadWriter interface {
	WriteTo(writer *WSPacketWriter)
	Length() uint64
}

const (
	MessagePayloadType = 0x00
	FilePayloadType    = 0x01
)
