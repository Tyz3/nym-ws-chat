package request

import (
	"fmt"
	. "nym-ws-chat/client/chat_payload"
	. "nym-ws-chat/client/web_socket_packet"
)

type SendRequest struct {
	request

	UsingSurb bool
	Address   string

	text string
	path string
}

func NewSendRequest(writer *WSPacketWriter, usingSurb bool, address string) *SendRequest {
	return &SendRequest{
		request: request{
			Tag:            SendRequestType,
			WSPacketWriter: writer,
		},
		UsingSurb: usingSurb,
		Address:   address,
	}
}

func (m *SendRequest) SetFile(path string) *SendRequest {
	m.path = path
	return m
}

func (m *SendRequest) SetMessage(text string) *SendRequest {
	m.text = text
	return m
}

func (m *SendRequest) Send() {
	m.WSPacketWriter.WriteByte(m.Tag) // 1 байт
	if m.UsingSurb {                  // 1 байт
		m.WSPacketWriter.WriteByte(0x01)
	} else {
		m.WSPacketWriter.WriteByte(0x00)
	}

	m.WSPacketWriter.WriteNymAddress(m.Address) // 96 байт

	if m.path != "" {
		var payload PayloadWriter = NewFilePayloadW(m.path)
		m.WSPacketWriter.WriteUint64(payload.Length()) // 8 байт
		payload.WriteTo(m.WSPacketWriter)              // M байт
	} else {
		var payload PayloadWriter = NewMessagePayloadW(m.text)
		fmt.Println("Payload length:", payload.Length())
		m.WSPacketWriter.WriteUint64(payload.Length()) // 8 байт
		payload.WriteTo(m.WSPacketWriter)              // M байт
	}
}
