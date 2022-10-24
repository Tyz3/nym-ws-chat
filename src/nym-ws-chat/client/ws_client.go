package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"nym-ws-chat/client/response"
	. "nym-ws-chat/client/web_socket_packet"
	"os"
)

type Client struct {
	host string
	port int
	conn *websocket.Conn

	Closed    bool
	Benchmark struct {
		N int
	}
}

func NewClient(host string, port int) *Client {
	client := &Client{
		host: host,
		port: port,
	}

	conn, resp, err := websocket.DefaultDialer.Dial(client.GetUrl(), nil)

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%s -> connection established (%d)\n", client.GetUrl(), resp.StatusCode)
	}

	client.conn = conn

	return client
}

func (c *Client) GetUrl() string {
	return fmt.Sprintf("ws://%s:%d", c.host, c.port)
}

func (c *Client) Close() {
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println("Error during closing websocket:", err)
		os.Exit(1)
	}
	c.Closed = true
}

func (c *Client) GetBinaryWriter() *WSPacketWriter {
	writer, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		panic(err)
	}
	return NewWSPacketWriter(websocket.BinaryMessage, writer)
}

func (c *Client) ReadSocketLoop() {
	for !c.Closed {
		messageType, reader, err := c.conn.NextReader()
		if err != nil {
			fmt.Println(err)
			c.Benchmark.N--
			continue
		}

		if messageType == -1 {
			fmt.Println("WebSocket closed externaly")
			c.Benchmark.N--
			return
		}

		// Создаём пакет из соединения
		packet := NewWSPacketReader(messageType, reader)
		if packet == nil {
			fmt.Println("WSPacketReader is nil")
			c.Benchmark.N--
			continue
		}

		// Печать пакета
		fmt.Println(packet.String())

		// Проверка нужного типа сообщения (TextMessage или BinaryMessage)
		if !packet.IsValid() {
			fmt.Println("Принятый тип сообщения не поддерживается, сообщение отброшено")
			c.Benchmark.N--
			continue
		}

		// Чтение пакета по типу
		switch packet.Type {
		case websocket.BinaryMessage:
			// Читаем первый байт пакета (сигнатура сообщения)
			sig := packet.ReadByte()

			resp := response.CreateResponse(sig, packet)
			if resp == nil {
				fmt.Printf("Тип сообщения 0x%02x не распознан\n", sig)
				continue
			}

			resp.Parse()
			fmt.Println(resp.String())
		case websocket.TextMessage:
			fmt.Println("Поддержка TextMessage не реализована")
		}
		c.Benchmark.N--
	}
}
