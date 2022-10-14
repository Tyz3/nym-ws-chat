package response

import (
	"io"
	"os"
	"strings"
)

type ErrorResponse struct {
	response

	Length  uint64
	Message string
}

func NewErrorResponse(reader io.Reader) *ErrorResponse {
	return &ErrorResponse{
		response: response{
			tag:    ErrorResponseType,
			reader: reader,
		},
	}
}

func (r *ErrorResponse) Parse() {
	file, _ := os.OpenFile("error", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)

	io.Copy(file, r.reader)
	file.Close()

	//// Читаем 8 байт - размер сообщения
	//msgLength := make([]byte, 8)
	//read, err := r.reader.Read(msgLength)
	//print("Read:", read)
	//if err != nil {
	//	panic(err)
	//}
	//r.Length = binary.BigEndian.Uint64(msgLength)
	//
	//// Читаем сообщение
	//message := bytes.NewBuffer(make([]byte, 0, r.Length))
	//written, err := io.CopyN(message, r.reader, int64(r.Length))
	//print("Written:", written)
	//if err != nil {
	//	panic(err)
	//}
	//r.Message = string(message.Bytes())
}

func (r *ErrorResponse) ToString() string {
	var sb strings.Builder
	sb.WriteString("Check file: error")
	//sb.WriteString("Tag: ")
	//sb.WriteString(fmt.Sprintf("0x%02x", r.tag))
	//sb.WriteString("\n")
	//sb.WriteString("Length: ")
	//sb.WriteString(fmt.Sprintf("%d", r.Length))
	//sb.WriteString("\n")
	//sb.WriteString("Message: ")
	//sb.WriteString(r.Message)
	return sb.String()
}
