package response

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"nym-ws-chat/client/chat_payload"
	"nym-ws-chat/client/web_socket_packet"
	"os"
	"strings"
)

type ReceiveResponse struct {
	response

	HasSurb    bool
	SurbLength uint64
	Surb       []byte

	payloadLength uint64
	Payload       chat_payload.Payload
}

func NewReceiveResponse(wsPacketReader *web_socket_packet.WSPacketReader) *ReceiveResponse {
	return &ReceiveResponse{
		response: response{
			Tag:            ReceiveResponseType,
			WSPacketReader: wsPacketReader,
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

	// Создаём объект полезной нагрузки
	// |                     payload(размером M)                  |
	// Message
	// |   sig   |                      text                      |
	// |  1 байт |                    M-1 байт                    |
	// File
	// |   sig   | file-name-length(N) | file-name | file-content |
	// |  1 байт |         8 байт      |   N байт  | M-1-8-N байт |
	r.Payload = chat_payload.CreatePayload(payloadSig, r.payloadLength, r.WSPacketReader)

	if r.Payload == nil {
		fmt.Fprintf(os.Stderr, "Тип (0x%0x) полезной нагрузки не распознан", payloadSig)
		return
	}

	r.Payload.Process()
}

func (r *ReceiveResponse) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("HasSurb: %+v\n", r.HasSurb))
	if r.HasSurb {
		sb.WriteString(fmt.Sprintf("SurbLength: %d\n", r.SurbLength))
		sb.WriteString(fmt.Sprintf("Surb: %s\n", base58.Encode(r.Surb)))
	}
	sb.WriteString(fmt.Sprintf("PayloadLength: %d\n", r.payloadLength))
	sb.WriteString(fmt.Sprintf("Payload: %s", r.Payload.String()))
	return sb.String()
}
