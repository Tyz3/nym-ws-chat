package request

import (
	"github.com/btcsuite/btcd/btcutil/base58"
	. "nym-ws-chat/client/chat_payload"
	. "nym-ws-chat/client/web_socket_packet"
)

type ReplyRequest struct {
	request

	Surb string

	text string
	path string
}

func NewReplyRequest(writer *WSPacketWriter, surb string) *ReplyRequest {
	return &ReplyRequest{
		request: request{
			Tag:            ReplyRequestType,
			WSPacketWriter: writer,
		},
		Surb: surb,
	}
}

func (m *ReplyRequest) SetFile(path string) *ReplyRequest {
	m.path = path
	return m
}

func (m *ReplyRequest) SetMessage(text string) *ReplyRequest {
	m.text = text
	return m
}

func (m *ReplyRequest) Send() {
	m.WSPacketWriter.WriteByte(m.Tag) // 1 байт

	surbBytes := base58.Decode(m.Surb)
	m.WSPacketWriter.WriteUint64(uint64(len(surbBytes))) // 8 байт
	m.WSPacketWriter.Write(surbBytes)                    // L байт

	if m.path != "" {
		var payload PayloadWriter = NewFilePayloadW(m.path)
		m.WSPacketWriter.WriteUint64(payload.Length()) // 8 байт
		payload.WriteTo(m.WSPacketWriter)              // M байт
	} else {
		var payload PayloadWriter = NewMessagePayloadW(m.text)
		m.WSPacketWriter.WriteUint64(payload.Length()) // 8 байт
		payload.WriteTo(m.WSPacketWriter)              // M байт
	}
}
