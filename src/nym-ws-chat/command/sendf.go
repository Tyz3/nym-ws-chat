package command

import (
	"fmt"
	. "nym-ws-chat/client"
	"nym-ws-chat/client/request"
	"nym-ws-chat/config"
	"os"
	"strconv"
)

type SendFCmd struct {
	command
}

func NewSendFCmd(name string, argsRequired int) *SendFCmd {
	return &SendFCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *SendFCmd) Execute(config *config.Config, args []string) {
	client := NewClient(config.Client.Host, config.Client.Port)

	contactNum, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Выбор контакта из списка
	if contactNum >= len(config.Contacts) {
		fmt.Println("Не найден контакт под номером", contactNum)
		return
	}
	contact := config.Contacts[contactNum]

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
	request.NewSendRequest(writer, true, contact.Address).SetFile(args[3]).Send()
	writer.Close()

	//client.Close()
	//cmd.command.done = true
}

func (cmd *SendFCmd) GetParams() string {
	return "<контакт> <путь_к_файлу>"
}

func (cmd *SendFCmd) GetDescription() string {
	return "отправить файл через сеть Nym"
}
