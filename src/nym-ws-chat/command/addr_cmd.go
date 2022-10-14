package command

import (
	"github.com/gorilla/websocket"
	. "nym-ws-chat/client"
	"nym-ws-chat/client/request"
	"nym-ws-chat/config"
)

type AddrCmd struct {
	command
}

func NewAddrCmd(name string, argsRequired int) *AddrCmd {
	return &AddrCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *AddrCmd) Execute(config *config.Config, args []string) {
	client := NewClient(config.Client.Host, config.Client.Port)

	go client.ReadSocket()

	writer, err := client.Conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		panic(err)
	}

	request.NewSelfAddressRequest().Send(writer)
	writer.Close()
}

func (cmd *AddrCmd) GetParams() string {
	return ""
}

func (cmd *AddrCmd) GetDescription() string {
	return "узнать собственный адрес в сети Nym"
}
