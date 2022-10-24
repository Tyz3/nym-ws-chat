package web_socket_packet

import (
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"nym-ws-chat/client/nym_util"
	"os"
	"strconv"
	"strings"
)

type WSPacketReader struct {
	Type   int
	reader io.Reader

	CurrentPacket []byte
}

func NewWSPacketReader(msgType int, reader io.Reader) *WSPacketReader {
	return &WSPacketReader{
		Type:   msgType,
		reader: reader,
	}
}

func (p *WSPacketReader) Read(buf []byte) {
	n, err := p.reader.Read(buf)
	if err != nil {
		if n != len(buf) {
			fmt.Fprintf(os.Stderr, "Размер буфера %d байт, но прочитано всего %d байт\n", len(buf), n)
		}
		panic(err)
	}
	p.CurrentPacket = append(p.CurrentPacket, buf...)
}

func (p *WSPacketReader) ReadByte() byte {
	buf := make([]byte, 1)
	p.Read(buf)

	return buf[0]
}

func (p *WSPacketReader) ReadUint64() uint64 {
	buf := make([]byte, 8)
	p.Read(buf)

	return binary.BigEndian.Uint64(buf)
}

func (p *WSPacketReader) ReadUint16() uint16 {
	buf := make([]byte, 2)
	p.Read(buf)

	return binary.BigEndian.Uint16(buf)
}

func (p *WSPacketReader) ReadN(n uint64) []byte {
	buf := make([]byte, n)
	p.Read(buf)

	return buf
}

func (p *WSPacketReader) ReadFlag() bool {
	flag := p.ReadByte()
	return flag == 0x01
}

func (p *WSPacketReader) ReadString(length uint64) string {
	buf := make([]byte, length)
	p.Read(buf)

	return string(buf)
}

func (p *WSPacketReader) ReadNymAddress() string {
	buf := make([]byte, 96)
	p.Read(buf)
	return nym_util.NymAddressFromBytes(buf)
}

func (p *WSPacketReader) Reader() io.Reader {
	return p.reader
}

func (p *WSPacketReader) IsValid() bool {
	return p.Type == websocket.BinaryMessage || p.Type == websocket.TextMessage
}

func (p *WSPacketReader) String() string {
	var sb strings.Builder

	sb.WriteString("Message Type: ")
	sb.WriteString(strconv.Itoa(p.Type))
	switch p.Type {
	case websocket.TextMessage:
		sb.WriteString(" (TextMessage)")
	case websocket.BinaryMessage:
		sb.WriteString(" (BinaryMessage)")
	}

	return sb.String()
}
