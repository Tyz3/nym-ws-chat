package response

import (
	"fmt"
	"nym-ws-chat/client/request"
	"nym-ws-chat/client/web_socket_packet"
)

type SelfAddressResponse struct {
	response

	Address []byte
}

func NewSelfAddressResponse(wsPacketReader *web_socket_packet.WSPacketReader) *SelfAddressResponse {
	return &SelfAddressResponse{
		response: response{
			Tag:            SelfAddressResponseType,
			WSPacketReader: wsPacketReader,
		},
	}
}

func (r *SelfAddressResponse) Parse() {
	// Читаем surb-байт
	r.Address = r.WSPacketReader.ReadN(96)
}

func (r *SelfAddressResponse) String() string {
	return fmt.Sprintf("Address: %s", request.NymAddressFromBytes(r.Address))
}
