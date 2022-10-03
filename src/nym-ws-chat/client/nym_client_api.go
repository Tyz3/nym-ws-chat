package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"nym-ws-chat/message"
	"os"
)

type Client struct {
	Host string
	Port int
	Conn *websocket.Conn

	Closed bool

	Benchmark struct {
		N int
	}
}

func NewClient(host string, port int) *Client {
	client := &Client{
		Host: host,
		Port: port,
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
	return fmt.Sprintf("ws://%s:%d", c.Host, c.Port)
}

func (c *Client) Close() {
	fmt.Println("Closing websocket...")
	err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println("Error during closing websocket:", err)
		os.Exit(1)
	}
	c.Closed = true
}

func (c *Client) SendMessage(message message.Message) {
	err := c.Conn.WriteMessage(websocket.TextMessage, message.ToJson())
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("[%s] %s\n", time.Now().Format(time.ANSIC), string(message.ToJson()))
}

func (c *Client) ReadSocket(outputChannel chan<- string) {
	for !c.Closed {
		_, data, err := c.Conn.ReadMessage()
		c.Benchmark.N--
		if err != nil {
			fmt.Println(err)
		} else {
			outputChannel <- string(data)
		}
	}
}

func (c *Client) StartPrint(inputChannel <-chan string) {
	for !c.Closed {
		<-inputChannel
		//msg := <-inputChannel
		//fmt.Printf("[%s] %s\n", time.Now().Format(time.ANSIC), msg)
	}
}

func (c *Client) PrintInputMessages() {
	channel := make(chan string, 10) // Канал для пересылки сообщений между горутинами
	go c.ReadSocket(channel)
	go c.StartPrint(channel)
}
