package response

import (
	"fmt"
	"io"
	"nym-ws-chat/client/request"
	"strings"
)

type SelfAddressResponse struct {
	response

	Address string
}

func NewSelfAddressResponse(reader io.Reader) *SelfAddressResponse {
	return &SelfAddressResponse{
		response: response{
			tag:    SelfAddressResponseType,
			reader: reader,
		},
	}
}

func (r *SelfAddressResponse) Parse() {
	// Читаем surb-байт
	address := make([]byte, 96)
	_, _ = r.reader.Read(address)
	r.Address = request.NymAddressFromBytes(address)
}

func (r *SelfAddressResponse) ToString() string {
	var sb strings.Builder
	sb.WriteString("Tag: ")
	sb.WriteString(fmt.Sprintf("0x%02x", r.tag))
	sb.WriteString("\n")
	sb.WriteString("Address: ")
	sb.WriteString(r.Address)
	return sb.String()
}
