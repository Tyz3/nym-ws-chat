package command

import (
	"fmt"
	. "nym-ws-chat/client"
	"nym-ws-chat/config"
	"nym-ws-chat/message"
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

	// Оnправка сообщения
	msg := message.NewOneWayMessage(text, contact.Address, false)
	client.SendMessage(msg)

	client.Close()

	cmd.command.done = true
}

func (cmd *SendCmd) GetParams() string {
	return "<контакт> <сообщение>"
}

func (cmd *SendCmd) GetDescription() string {
	return "отправить сообщение через сеть Nym"
}
