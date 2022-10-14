package chat_payload

import (
	"fmt"
	"io"
	"nym-ws-chat/client/web_socket_packet"
	"os"
)

type File struct {
	payload
	FileName   string
	FileLength int64
}

func NewFile(length uint64, wsPacketReader *web_socket_packet.WSPacketReader) *File {
	return &File{
		payload: payload{
			Length:         length,
			Sig:            FilePayloadType,
			WSPacketReader: wsPacketReader,
		},
	}
}

func (p *File) Process() {

	// Длина имени файла
	fileNameLength := p.WSPacketReader.ReadUint64()

	// Имя файла
	p.FileName = p.WSPacketReader.ReadString(fileNameLength)

	// Создаём файл
	file, err := os.OpenFile(p.FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось создать файл с именем %s, err = %s\n", p.FileName, err.Error())
		return
	}

	// Копируем данные файла из потока в файл
	p.FileLength, err = io.Copy(file, p.WSPacketReader.Reader())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось записать данные в файл %s, written = %d, err = %s\n", p.FileName, p.FileLength, err.Error())
		return
	}

	file.Close()
}

func (p *File) String() string {
	return fmt.Sprintf("Sig: 0x%02x, File: %s (%d байт)", p.Sig, p.FileName, p.FileLength)
}
