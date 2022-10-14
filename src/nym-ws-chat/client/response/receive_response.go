package response

import (
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"io"
	"nym-ws-chat/client/request"
	"os"
	"strings"
)

type ReceiveResponse struct {
	response

	HasSurb    bool
	SurbLength uint64
	Surb       []byte

	PayloadLength uint64
	Payload       string

	IsFile         bool
	FileNameLength uint64
	FileName       string
}

func NewReceiveResponse(reader io.Reader) *ReceiveResponse {
	return &ReceiveResponse{
		response: response{
			tag:    ReceiveResponseType,
			reader: reader,
		},
	}
}

func (r *ReceiveResponse) Parse() {
	// Читаем surb-байт
	hasSurb := make([]byte, 1)
	_, _ = r.reader.Read(hasSurb)
	if hasSurb[0] == 0x01 {
		r.HasSurb = true
	}

	if r.HasSurb {
		// Читаем длину surb
		surbLength := make([]byte, 8)
		_, _ = r.reader.Read(surbLength)
		r.SurbLength = binary.BigEndian.Uint64(surbLength)

		// Читаем surb
		surb := make([]byte, r.SurbLength)
		_, _ = r.reader.Read(surb)
		r.Surb = surb
	}

	// Читаем длину полезной нагрузки
	payloadLength := make([]byte, 8)
	_, _ = r.reader.Read(payloadLength)
	r.PayloadLength = binary.BigEndian.Uint64(payloadLength)

	// Читаем полезную нагрузку
	isFileByte := make([]byte, 1)
	_, _ = r.reader.Read(isFileByte)

	// Если 0x00, то сообщение - файл, если 0x01 - текст
	if isFileByte[0] == request.MesgByte[0] {
		payload := make([]byte, r.PayloadLength-1)
		_, _ = r.reader.Read(payload)
		r.Payload = string(payload)
	} else if isFileByte[0] == request.FileByte[0] {
		r.IsFile = true

		// Длина имени файла
		fileNameLength := make([]byte, 8)
		r.reader.Read(fileNameLength)
		r.FileNameLength = binary.BigEndian.Uint64(fileNameLength)

		// Имя файла
		fileName := make([]byte, r.FileNameLength)
		r.reader.Read(fileName)
		r.FileName = string(fileName)

		// Создаём файл
		file, err := os.OpenFile(r.FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}

		// Копируем данные файла из потока в файл
		written, err := io.Copy(file, r.reader)
		if err != nil {
			panic(err)
		}

		file.Close()

		fmt.Printf("Сохранён файл '%s', размер: %d байт\n", fileName, written)
	} else {
		fmt.Println("Принято сообщение с неизвестным содержимым")
	}
}

func (r *ReceiveResponse) ToString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tag: 0x%02x\n", r.tag))
	sb.WriteString(fmt.Sprintf("HasSurb: %+v\n", r.HasSurb))
	if r.HasSurb {
		sb.WriteString(fmt.Sprintf("SurbLength: %d\n", r.SurbLength))
		sb.WriteString(fmt.Sprintf("Surb: %s\n", base58.Encode(r.Surb)))
	}
	sb.WriteString(fmt.Sprintf("PayloadLength: %d\n", r.PayloadLength))
	if !r.IsFile {
		sb.WriteString(fmt.Sprintf("Payload: %s", r.Payload))
	} else {
		sb.WriteString("Payload: file")
	}
	return sb.String()
}
