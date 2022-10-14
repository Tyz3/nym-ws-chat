package response

import (
	"io"
)

type response struct {
	tag    byte
	reader io.Reader
}

type Response interface {
	ToString() string
	Parse()
}

const (
	ErrorResponseType       = 0x00
	ReceiveResponseType     = 0x01
	SelfAddressResponseType = 0x02
)

func CreateResponse(reader io.Reader) (Response, byte) {
	tag := make([]byte, 1)
	_, err := reader.Read(tag)
	if err != nil {
		return nil, tag[0]
	}

	switch tag[0] {
	case ErrorResponseType:
		return NewErrorResponse(reader), tag[0]
	case ReceiveResponseType:
		return NewReceiveResponse(reader), tag[0]
	case SelfAddressResponseType:
		return NewSelfAddressResponse(reader), tag[0]
	default:
		return nil, tag[0]
	}
}
