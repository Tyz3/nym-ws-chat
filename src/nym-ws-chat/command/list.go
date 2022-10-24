package command

import (
	"fmt"
	"nym-ws-chat/config"
	"strings"
)

type ListCmd struct {
	command
}

func NewListCmd(name string, argsRequired int) *ListCmd {
	return &ListCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *ListCmd) Execute(cfg *config.Config, args []string) {
	var sb strings.Builder

	sb.WriteString("Список контактов:\n")
	for i, contact := range cfg.Contacts {
		sb.WriteString(fmt.Sprintf("#%d %-8s %s\n", i, contact.Alias, contact.Address))
	}

	fmt.Println(sb.String())

	cmd.command.done = true
}

func (cmd *ListCmd) GetParams() string {
	return ""
}

func (cmd *ListCmd) GetDescription() string {
	return "отобразить список контактов"
}
