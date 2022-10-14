package request

import (
	"io"
)

type SelfAddressRequest struct {
	request
}

func NewSelfAddressRequest() *SelfAddressRequest {
	return &SelfAddressRequest{
		request{
			tag: 0x02,
		},
	}
}

func (m *SelfAddressRequest) Send(writer io.WriteCloser) {
	_, _ = writer.Write([]byte{m.tag})
}
