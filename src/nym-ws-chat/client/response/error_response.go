package response

import (
	"fmt"
	"io"
	. "nym-ws-chat/client/web_socket_packet"
	"os"
)

type ErrorResponse struct {
	response

	payloadLength uint64
}

func NewErrorResponse(reader *WSPacketReader) *ErrorResponse {
	return &ErrorResponse{
		response: response{
			Tag:            ErrorResponseType,
			WSPacketReader: reader,
		},
	}
}

func (r *ErrorResponse) Parse() {
	file, _ := os.OpenFile("error", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)

	r.WSPacketReader.ReadN(9)

	io.Copy(file, r.WSPacketReader.Reader())
	file.Close()
}

func (r *ErrorResponse) String() string {
	return fmt.Sprintf("Error. Check the file.")
}
