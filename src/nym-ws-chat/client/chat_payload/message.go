package chat_payload

import (
	"fmt"
	. "nym-ws-chat/client/web_socket_packet"
)

type MessagePayload struct {
	payload

	Text string
}

func NewMessagePayloadR(length uint64) *MessagePayload {
	return &MessagePayload{
		payload: payload{
			Sig:    MessagePayloadType,
			Length: length,
		},
	}
}

func NewMessagePayloadW(text string) *MessagePayload {
	return &MessagePayload{
		payload: payload{
			Sig:    MessagePayloadType,
			Length: uint64(len(text) + 1),
		},
		Text: text,
	}
}

func (p *MessagePayload) ReadFrom(reader *WSPacketReader) {
	p.Text = reader.ReadString(p.payload.Length - 1)
}

func (p *MessagePayload) String() string {
	return fmt.Sprintf("Sig: 0x%02x, Text: %s", p.Sig, p.Text)
}

func (p *MessagePayload) WriteTo(writer *WSPacketWriter) {
	writer.WriteByte(p.Sig)
	writer.WriteString(p.Text)
}

func (p *MessagePayload) Length() uint64 {
	return p.payload.Length
}
