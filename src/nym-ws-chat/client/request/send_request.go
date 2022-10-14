package request

import (
	"encoding/binary"
	"fmt"
	"io"
	"nym-ws-chat/client/chat_payload"
	"os"
)

type SendRequest struct {
	request

	surb    byte
	address [96]byte

	message       string
	file          *os.File
	fileInfo      os.FileInfo
	payloadLength []byte
}

func NewSendRequest(surb bool, address string) *SendRequest {
	r := &SendRequest{
		request: request{
			tag: 0x00,
		},
		address:       NymAddressToBytes(address),
		payloadLength: make([]byte, 8),
	}

	fmt.Println("Адрес для отправки сообщения:\n", NymAddressFromBytes(r.address[:]))

	if surb {
		r.surb = 0x01
	} else {
		r.surb = 0x00
	}

	return r
}

func (m *SendRequest) SetFile(fileInfo os.FileInfo, path string) *SendRequest {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия файла: %s (Size:%d, Mode:%s, IsDir:%+v)\n", path, fileInfo.Size(), fileInfo.Mode().String(), fileInfo.IsDir())
		panic(err)
	}

	m.file = file
	m.fileInfo = fileInfo
	binary.BigEndian.PutUint64(m.payloadLength, uint64(fileInfo.Size()+8+1+int64(len(fileInfo.Name()))))
	return m
}

func (m *SendRequest) SetMessage(text string) *SendRequest {
	m.message = text
	binary.BigEndian.PutUint64(m.payloadLength, uint64(len(text)+1))
	return m
}

func (m *SendRequest) Send(writer io.WriteCloser) {
	writer.Write([]byte{m.tag, m.surb})
	writer.Write(m.address[:])
	writer.Write(m.payloadLength)

	if m.file != nil {
		// Сигнатура файла
		writer.Write([]byte{chat_payload.FilePayloadType})

		// Имя файла
		fileName := m.fileInfo.Name()
		fileNameLength := make([]byte, 8)
		binary.BigEndian.PutUint64(fileNameLength, uint64(len(fileName)))
		writer.Write(fileNameLength)
		writer.Write([]byte(fileName))

		// Содержимое файла
		written, err := io.Copy(writer, m.file)
		m.file.Close()
		fmt.Println("File Written:", written)
		if err != nil {
			panic(err)
		}
	} else {
		// Сигнатура текста
		writer.Write([]byte{chat_payload.MessagePayloadType})
		written, err := writer.Write([]byte(m.message))
		fmt.Println("Message Written:", written)
		if err != nil {
			panic(err)
		}
	}
}
