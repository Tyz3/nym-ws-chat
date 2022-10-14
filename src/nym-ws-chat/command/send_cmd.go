package command

import (
	"fmt"
	"github.com/gorilla/websocket"
	. "nym-ws-chat/client"
	"nym-ws-chat/client/request"
	"nym-ws-chat/config"
	"strconv"
	"strings"
)

type SendCmd struct {
	command
}

func NewSendCmd(name string, argsRequired int) *SendCmd {
	return &SendCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *SendCmd) Execute(config *config.Config, args []string) {
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
	text := strings.Join(args[3:], " ")

	// Включаем чтение сокета
	//go client.ReadSocket()

	// Отправка сообщения
	writer, err := client.Conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		panic(err)
	}
	request.NewSendRequest(true, contact.Address).SetMessage(text).Send(writer)
	writer.Close()

	client.Close()
	cmd.command.done = true
}

func (cmd *SendCmd) GetParams() string {
	return "<контакт> <сообщение>"
}

func (cmd *SendCmd) GetDescription() string {
	return "отправить сообщение через сеть Nym"
}
