package request

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type SendRequest struct {
	request

	surb    byte
	address [96]byte

	message       string
	file          *os.File
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

	fmt.Println("Адрес для отправки сообщения:\n", r.address[:], "\n", NymAddressFromBytes(r.address[:]))

	if surb {
		r.surb = 0x01
	} else {
		r.surb = 0x00
	}

	return r
}

func (m *SendRequest) SetFile(fileInfo os.FileInfo) *SendRequest {
	file, err := os.Open(fileInfo.Name())
	if err != nil {
		panic(err)
	}

	m.file = file
	binary.BigEndian.PutUint64(m.payloadLength, uint64(fileInfo.Size()+8+1+int64(len(fileInfo.Name()))))
	return m
}

func (m *SendRequest) SetMessage(text string) *SendRequest {
	m.message = text
	binary.BigEndian.PutUint64(m.payloadLength, uint64(len(text)+1))
	return m
}

func (m *SendRequest) Send(writer io.WriteCloser) {
	_, _ = writer.Write([]byte{m.tag, m.surb})
	_, _ = writer.Write(m.address[:])
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
