package chat_payload

import (
	"fmt"
	"io"
	. "nym-ws-chat/client/web_socket_packet"
	"os"
)

type FilePayload struct {
	payload

	Info os.FileInfo
	Path string
}

func NewFilePayloadR(length uint64) *FilePayload {
	return &FilePayload{
		payload: payload{
			Sig:    FilePayloadType,
			Length: length,
		},
	}
}

func NewFilePayloadW(path string) *FilePayload {
	p := &FilePayload{
		payload: payload{
			Sig: FilePayloadType,
		},
		Path: path,
	}

	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	p.Info = info
	p.payload.Length = uint64(info.Size()) + uint64(len(info.Name())) + 1 + 2

	return p
}

func (p *FilePayload) ReadFrom(reader *WSPacketReader) {
	// Длина имени файла
	fileNameLength := reader.ReadUint64()

	// Имя файла
	fileName := reader.ReadString(fileNameLength)

	// Создаём файл
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось создать файл с именем %s, err = %s\n", fileName, err.Error())
		return
	}
	p.Path = file.Name()

	// Копируем данные из потока в файл
	written, err := io.CopyN(file, reader.Reader(), int64(p.payload.Length-1-2-fileNameLength))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось записать данные в файл %s, written = %d, err = %s\n", fileName, written, err.Error())
		return
	}

	info, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось получить информацию по файлу %s, err = %s\n", fileName, err.Error())
		return
	}
	file.Close()

	p.Info = info
}

func (p *FilePayload) String() string {
	return fmt.Sprintf("Sig: 0x%02x, File: %s (%d байт)", p.Sig, p.Info.Name(), p.Info.Size())
}

func (p *FilePayload) WriteTo(writer *WSPacketWriter) {
	fileName := p.Info.Name()

	writer.WriteByte(p.Sig)
	writer.WriteUint16(uint16(len(fileName))) // 2 байта
	writer.WriteString(fileName)              // N байт

	file, err := os.Open(p.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось открыть файл %s на чтение, err = %s\n", p.Path, err.Error())
		return
	}

	written, err := io.Copy(writer.Writer(), file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось скопировать файл %s в поток вывода, err = %s\n", file.Name(), err.Error())
		return
	}
	fmt.Println("Written:", written)
}

func (p *FilePayload) Length() uint64 {
	return p.payload.Length
}
