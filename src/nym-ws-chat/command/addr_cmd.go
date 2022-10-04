package command

import (
	. "nym-ws-chat/client"
	"nym-ws-chat/config"
	"nym-ws-chat/message"
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

	channel := make(chan string, 10) // Канал для пересылки сообщений между горутинами
	go client.ReadSocket(channel)
	go client.StartPrint(channel)

	msg := message.NewSelfAddressMessage()
	client.SendMessage(msg)
}

func (cmd *AddrCmd) GetParams() string {
	return ""
}

func (cmd *AddrCmd) GetDescription() string {
	return "узнать собственный адрес в сети Nym"
}
