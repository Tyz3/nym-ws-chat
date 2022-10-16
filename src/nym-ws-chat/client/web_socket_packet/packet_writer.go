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

type WSPacketWriter struct {
	Type   int
	writer io.WriteCloser
}

func NewWSPacketWriter(msgType int, writer io.WriteCloser) *WSPacketWriter {
	return &WSPacketWriter{
		Type:   msgType,
		writer: writer,
	}
}

func (p *WSPacketWriter) Write(buf []byte) {
	n, err := p.writer.Write(buf)
	if err != nil {
		if n != len(buf) {
			fmt.Fprintf(os.Stderr, "Размер буфера %d байт, но записано всего %d байт\n", len(buf), n)
		}
		panic(err)
	}
}

func (p *WSPacketWriter) WriteByte(b ...byte) {
	p.Write(b)
}

func (p *WSPacketWriter) WriteUint64(num uint64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, num)
	p.Write(buf)
}

func (p *WSPacketWriter) WriteFlag(flag bool) {
	if flag {
		p.Write([]byte{0x01})
	} else {
		p.Write([]byte{0x00})
	}
}

func (p *WSPacketWriter) WriteString(s string) {
	p.Write([]byte(s))
}

func (p *WSPacketWriter) WriteNymAddress(address string) {
	fmt.Println("Address:", address)
	addrBytes := nym_util.NymAddressToBytes(address)
	fmt.Println("Address:", addrBytes[:])
	p.Write(addrBytes[:])
}

func (p *WSPacketWriter) Writer() io.WriteCloser {
	return p.writer
}

func (p *WSPacketWriter) Close() {
	p.writer.Close()
}

func (p *WSPacketWriter) String() string {
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
