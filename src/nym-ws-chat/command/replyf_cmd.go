package command

import (
	"github.com/gorilla/websocket"
	. "nym-ws-chat/client"
	"nym-ws-chat/client/request"
	"nym-ws-chat/config"
	"os"
)

type ReplyFCmd struct {
	command
}

func NewReplyFCmd(name string, argsRequired int) *ReplyFCmd {
	return &ReplyFCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *ReplyFCmd) Execute(config *config.Config, args []string) {
	client := NewClient(config.Client.Host, config.Client.Port)

	surbBase58 := args[2]

	fileInfo, err := os.Stat(args[3])
	if err != nil {
		panic(err)
	}

	// Отправка сообщения
	writer, err := client.Conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		panic(err)
	}
	request.NewReplyRequest(surbBase58).SetFile(fileInfo).Send(writer)
	writer.Close()

	client.Close()
	cmd.command.done = true
}

func (cmd *ReplyFCmd) GetParams() string {
	return "<surb> <путь_к_файлу>"
}

func (cmd *ReplyFCmd) GetDescription() string {
	return "ответить файлом на сообщение"
}
