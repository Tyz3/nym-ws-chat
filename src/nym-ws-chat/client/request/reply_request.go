package request

import (
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"io"
	"os"
)

type ReplyRequest struct {
	request

	surbLength []byte
	surb       []byte

	payloadLength []byte
	message       string
	file          *os.File
}

func NewReplyRequest(surb string) *ReplyRequest {
	r := &ReplyRequest{
		request: request{
			tag: 0x01,
		},
		payloadLength: make([]byte, 8),
		surbLength:    make([]byte, 8),
		surb:          base58.Decode(surb),
	}

	binary.BigEndian.PutUint64(r.surbLength, uint64(len(r.surb)))

	return r
}

func (m *ReplyRequest) SetFile(fileInfo os.FileInfo) *ReplyRequest {
	file, err := os.Open(fileInfo.Name())
	if err != nil {
		panic(err)
	}

	m.file = file
	binary.BigEndian.PutUint64(m.payloadLength, uint64(fileInfo.Size()+8+1+int64(len(fileInfo.Name()))))
	return m
}

func (m *ReplyRequest) SetMessage(text string) *ReplyRequest {
	m.message = text
	binary.BigEndian.PutUint64(m.payloadLength, uint64(len(text)+1))
	return m
}

func (m *ReplyRequest) Send(writer io.WriteCloser) {
	_, _ = writer.Write([]byte{m.tag})
	_, _ = writer.Write(m.surbLength)
	_, _ = writer.Write(m.surb)
	_, _ = writer.Write(m.payloadLength)

	if m.file != nil {
		// Сигнатура файла
		_, _ = writer.Write(FileByte)

		// Имя файла
		fileName := m.file.Name()
		fileNameLength := make([]byte, 8)
		binary.BigEndian.PutUint64(fileNameLength, uint64(len(fileName)))
		_, _ = writer.Write(fileNameLength)
		_, _ = writer.Write([]byte(fileName))

		// Содержимое файла
		written, err := io.Copy(writer, m.file)
		m.file.Close()
		fmt.Println("File Written:", written)
		if err != nil {
			panic(err)
		}
	} else {
		// Сигнатура текста
		_, _ = writer.Write(MesgByte)
		written, err := writer.Write([]byte(m.message))
		fmt.Println("Message Written:", written)
		if err != nil {
			panic(err)
		}
	}
}
