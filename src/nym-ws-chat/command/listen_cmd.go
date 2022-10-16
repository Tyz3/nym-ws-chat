package command

import (
	. "nym-ws-chat/client"
	"nym-ws-chat/config"
)

type ListenCmd struct {
	command
}

func NewListenCmd(name string, argsRequired int) *ListenCmd {
	return &ListenCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *ListenCmd) Execute(config *config.Config, args []string) {
	client := NewClient(config.Client.Host, config.Client.Port)

	go client.ReadSocketLoop()
}

func (cmd *ListenCmd) GetParams() string {
	return ""
}

func (cmd *ListenCmd) GetDescription() string {
	return "вывести входящие сообщения из сети Nym"
}
