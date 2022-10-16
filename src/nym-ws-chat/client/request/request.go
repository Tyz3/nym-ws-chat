package request

import (
	"nym-ws-chat/client/web_socket_packet"
)

type request struct {
	Tag            byte
	WSPacketWriter *web_socket_packet.WSPacketWriter
}

type Request interface {
	Send()
}

const (
	SendRequestType        = 0x00
	ReplyRequestType       = 0x01
	SelfAddressRequestType = 0x02
)
