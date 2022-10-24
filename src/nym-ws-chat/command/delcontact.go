package command

import (
	"fmt"
	"nym-ws-chat/config"
	"strconv"
)

type DelContactCmd struct {
	command
}

func NewDelContactCmd(name string, argsRequired int) *DelContactCmd {
	return &DelContactCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *DelContactCmd) Execute(cfg *config.Config, args []string) {
	contactNum, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg.Contacts = append(cfg.Contacts[:contactNum], cfg.Contacts[contactNum+1:]...)
	cfg.Save()

	cmd.command.done = true
}

func (cmd *DelContactCmd) GetParams() string {
	return "<контакт>"
}

func (cmd *DelContactCmd) GetDescription() string {
	return "удалить контакт из списка контактов"
}
