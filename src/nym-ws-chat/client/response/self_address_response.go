package response

import (
	"fmt"
	. "nym-ws-chat/client/web_socket_packet"
)

type SelfAddressResponse struct {
	response

	Address string
}

func NewSelfAddressResponse(reader *WSPacketReader) *SelfAddressResponse {
	return &SelfAddressResponse{
		response: response{
			Tag:            SelfAddressResponseType,
			WSPacketReader: reader,
		},
	}
}

func (r *SelfAddressResponse) Parse() {
	// Читаем surb-байт
	r.Address = r.WSPacketReader.ReadNymAddress()
}

func (r *SelfAddressResponse) String() string {
	return fmt.Sprintf("Address: %s", r.Address)
}
