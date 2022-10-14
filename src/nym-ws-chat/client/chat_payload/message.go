package chat_payload

import (
	"fmt"
	"nym-ws-chat/client/web_socket_packet"
)

type Message struct {
	payload
	Text string
}

func NewMessage(length uint64, wsPacketReader *web_socket_packet.WSPacketReader) *Message {
	return &Message{
		payload: payload{
			Length:         length,
			Sig:            MessagePayloadType,
			WSPacketReader: wsPacketReader,
		},
	}
}

func (p *Message) Process() {
	p.Text = p.WSPacketReader.ReadString(p.Length - 1)
}

func (p *Message) String() string {
	return fmt.Sprintf("Sig: 0x%02x, Text: %s", p.Sig, p.Text)
}
