package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"nym-ws-chat/client/response"
	"os"
	"strings"
)

type Client struct {
	host string
	port int

	Conn      *websocket.Conn
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

	client.Conn = conn

	return client
}

func (c *Client) GetUrl() string {
	return fmt.Sprintf("ws://%s:%d", c.host, c.port)
}

func (c *Client) Close() {
	err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println("Error during closing websocket:", err)
		os.Exit(1)
	}
	c.Closed = true
}

//func (c *Client) SendMessage(message message.Message) {
//	err := c.Conn.WriteMessage(websocket.TextMessage, message.ToJson())
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//}

//func (c *Client) ReadSocket(outputChannel chan<- string) {
//	for !c.Closed {
//		_, data, err := c.Conn.ReadMessage()
//		c.Benchmark.N--
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			outputChannel <- string(data)
//		}
//	}
//}

func (c *Client) ReadSocket() {
	for !c.Closed {
		messageType, reader, err := c.Conn.NextReader()
		c.Benchmark.N--
		if err != nil {
			fmt.Println(err)
		}

		if messageType == -1 {
			return
		}

		var sb strings.Builder
		sb.WriteString("Received message type: ")
		if messageType == websocket.BinaryMessage {
			sb.WriteString("BinaryMessage")
		} else if messageType == websocket.TextMessage {
			sb.WriteString("TextMessage")
		}
		sb.WriteString(fmt.Sprintf("(%d)", messageType))
		fmt.Println(sb.String())

		if messageType == websocket.BinaryMessage {
			resp, tag := response.CreateResponse(reader)
			if resp == nil {
				fmt.Println("Тип сообщения", tag, "не распознан")
				continue
			}

			resp.Parse()
			fmt.Println(resp.ToString())

		} else if messageType == websocket.TextMessage {

		}
	}
}

//func (c *Client) StartPrint(inputChannel <-chan string) {
//	for !c.Closed {
//		msg := <-inputChannel
//
//		fmt.Printf("[%s] %s\n", time.Now().Format(time.ANSIC), msg)
//	}
//}
