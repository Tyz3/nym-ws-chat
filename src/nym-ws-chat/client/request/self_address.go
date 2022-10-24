package request

import (
	. "nym-ws-chat/client/web_socket_packet"
)

type SelfAddressRequest struct {
	request
}

func NewSelfAddressRequest(writer *WSPacketWriter) *SelfAddressRequest {
	return &SelfAddressRequest{
		request{
			Tag:            SelfAddressRequestType,
			WSPacketWriter: writer,
		},
	}
}

func (m *SelfAddressRequest) Send() {
	m.WSPacketWriter.WriteByte(m.Tag)
}
