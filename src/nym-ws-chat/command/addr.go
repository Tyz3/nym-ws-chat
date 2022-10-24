package command

import (
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

	go client.ReadSocketLoop()

	writer := client.GetBinaryWriter()
	request.NewSelfAddressRequest(writer).Send()
	writer.Close()
}

func (cmd *AddrCmd) GetParams() string {
	return ""
}

func (cmd *AddrCmd) GetDescription() string {
	return "узнать собственный адрес в сети Nym"
}
