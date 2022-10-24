package command

import (
	"fmt"
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

	if fileInfo.IsDir() {
		fmt.Println("Нужно указать файл, а не каталог")
		return
	}

	go client.ReadSocketLoop()

	// Отправка сообщения
	writer := client.GetBinaryWriter()
	request.NewReplyRequest(writer, surbBase58).SetFile(args[3]).Send()
	writer.Close()

	//client.Close()
	//cmd.command.done = true
}

func (cmd *ReplyFCmd) GetParams() string {
	return "<surb> <путь_к_файлу>"
}

func (cmd *ReplyFCmd) GetDescription() string {
	return "ответить файлом на сообщение"
}
