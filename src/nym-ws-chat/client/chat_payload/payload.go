package chat_payload

import (
	"nym-ws-chat/client/web_socket_packet"
)

type payload struct {
	Length         uint64
	Sig            byte
	WSPacketReader *web_socket_packet.WSPacketReader
}

type Payload interface {
	Process()
	String() string
}

const (
	MessagePayloadType = 0x00
	FilePayloadType    = 0x01
)

func CreatePayload(sig byte, length uint64, wsPacketReader *web_socket_packet.WSPacketReader) Payload {
	switch sig {
	case MessagePayloadType:
		return NewMessage(length, wsPacketReader)
	case FilePayloadType:
		return NewFile(length, wsPacketReader)
	default:
		return nil
	}
}
