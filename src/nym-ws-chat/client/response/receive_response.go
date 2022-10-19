package response

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	. "nym-ws-chat/client/chat_payload"
	. "nym-ws-chat/client/web_socket_packet"
	"os"
	"strings"
)

type ReceiveResponse struct {
	response

	HasSurb    bool
	SurbLength uint64
	Surb       []byte

	payloadLength uint64
	Payload       PayloadReader
}

func NewReceiveResponse(reader *WSPacketReader) *ReceiveResponse {
	return &ReceiveResponse{
		response: response{
			Tag:            ReceiveResponseType,
			WSPacketReader: reader,
		},
	}
}

func (r *ReceiveResponse) Parse() {
	// | surb-flag | surb-length(L) |   surb   | payload-length(M) |   payload   |
	// |   1 байт  |     8 байт     |  L байт  |       8 байт      |    M байт   |

	// Читаем surb-байт
	r.HasSurb = r.WSPacketReader.ReadFlag()

	if r.HasSurb {
		// Читаем длину surb
		r.SurbLength = r.WSPacketReader.ReadUint64()

		// Читаем surb
		r.Surb = r.WSPacketReader.ReadN(r.SurbLength)
	}

	// Читаем длину полезной нагрузки
	r.payloadLength = r.WSPacketReader.ReadUint64()

	// Читаем флаг полезной нагрузки
	payloadSig := r.WSPacketReader.ReadByte()

	if r.payloadLength > 400*1024*1024 { // 400 MiB
		fmt.Println("/Incorrect packet length:", r.payloadLength/1024/1024, "MiB (Max 400 MiB)")
		fmt.Println("\\Read Packet:", r.WSPacketReader.CurrentPacket)
		return
	}
	fmt.Println("Read Packet:", r.WSPacketReader.CurrentPacket)

	// Создаём объект полезной нагрузки
	// |                     payload(размером M)                  |
	// Message
	// |   sig   |                      text                      |
	// |  1 байт |                    M-1 байт                    |
	// File
	// |   sig   | file-name-length(N) | file-name | file-content |
	// |  1 байт |         8 байт      |   N байт  | M-1-8-N байт |
	switch payloadSig {
	case MessagePayloadType:
		fmt.Printf("NewMessagePayloadR(%d)\n", r.payloadLength)
		r.Payload = NewMessagePayloadR(r.payloadLength)
	case FilePayloadType:
		fmt.Printf("NewFilePayloadR(%d)\n", r.payloadLength)
		r.Payload = NewFilePayloadR(r.payloadLength)
	default:
		fmt.Fprintf(os.Stderr, "Тип (payloadSig=0x%02x) полезной нагрузки не распознан\n", payloadSig)
		return
	}

	r.Payload.ReadFrom(r.WSPacketReader)
}

func (r *ReceiveResponse) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("HasSurb: %+v\n", r.HasSurb))
	if r.HasSurb {
		sb.WriteString(fmt.Sprintf("SurbLength: %d\n", r.SurbLength))
		sb.WriteString(fmt.Sprintf("Surb: %s\n", base58.Encode(r.Surb)))
	}
	sb.WriteString(fmt.Sprintf("PayloadLength: %d\n", r.payloadLength))
	if r.Payload != nil {
		sb.WriteString(fmt.Sprintf("Payload: %s", r.Payload.String()))
	}
	return sb.String()
}
