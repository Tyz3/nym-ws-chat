package response

import (
	. "nym-ws-chat/client/web_socket_packet"
)

type response struct {
	Tag            byte
	WSPacketReader *WSPacketReader
}

type Response interface {
	String() string
	Parse()
}

const (
	ErrorResponseType       = 0x00
	ReceiveResponseType     = 0x01
	SelfAddressResponseType = 0x02
)

func CreateResponse(sig byte, wsPacketReader *WSPacketReader) Response {
	switch sig {
	case ErrorResponseType:
		return NewErrorResponse(wsPacketReader)
	case ReceiveResponseType:
		return NewReceiveResponse(wsPacketReader)
	case SelfAddressResponseType:
		return NewSelfAddressResponse(wsPacketReader)
	default:
		return nil
	}
}
