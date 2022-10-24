package command

import (
	"nym-ws-chat/config"
)

type AddContactCmd struct {
	command
}

func NewAddContactCmd(name string, argsRequired int) *AddContactCmd {
	return &AddContactCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *AddContactCmd) Execute(cfg *config.Config, args []string) {
	alias := args[2]
	address := args[3]
	cfg.Contacts = append(cfg.Contacts, config.Contact{Address: address, Alias: alias})
	cfg.Save()

	cmd.command.done = true
}

func (cmd *AddContactCmd) GetParams() string {
	return "<имя> <адрес>"
}

func (cmd *AddContactCmd) GetDescription() string {
	return "добавить новый контакт в список контактов"
}
