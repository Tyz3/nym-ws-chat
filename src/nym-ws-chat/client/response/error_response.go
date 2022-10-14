package response

import (
	"fmt"
	"io"
	"nym-ws-chat/client/web_socket_packet"
	"os"
)

type ErrorResponse struct {
	response

	payloadLength uint64
}

func NewErrorResponse(wsPacketReader *web_socket_packet.WSPacketReader) *ErrorResponse {
	return &ErrorResponse{
		response: response{
			Tag:            ErrorResponseType,
			WSPacketReader: wsPacketReader,
		},
	}
}

func (r *ErrorResponse) Parse() {
	file, _ := os.OpenFile("error", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)

	io.Copy(file, r.WSPacketReader.Reader())
	file.Close()
}

func (r *ErrorResponse) String() string {
	return fmt.Sprintf("Error. Check the file.")
}
