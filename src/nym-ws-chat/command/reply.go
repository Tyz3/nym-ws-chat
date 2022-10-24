package command

import (
	. "nym-ws-chat/client"
	"nym-ws-chat/client/request"
	"nym-ws-chat/config"
	"strings"
)

type ReplyCmd struct {
	command
}

func NewReplyCmd(name string, argsRequired int) *ReplyCmd {
	return &ReplyCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *ReplyCmd) Execute(config *config.Config, args []string) {
	client := NewClient(config.Client.Host, config.Client.Port)

	surbBase58 := args[2]

	text := strings.Join(args[3:], " ")

	// Включаем чтение сокета
	go client.ReadSocketLoop()

	// Отправка сообщения
	writer := client.GetBinaryWriter()
	request.NewReplyRequest(writer, surbBase58).SetMessage(text).Send()
	writer.Close()

	//client.Close()
	//cmd.command.done = true
}

func (cmd *ReplyCmd) GetParams() string {
	return "<surb> <сообщение>"
}

func (cmd *ReplyCmd) GetDescription() string {
	return "ответить на сообщение"
}
