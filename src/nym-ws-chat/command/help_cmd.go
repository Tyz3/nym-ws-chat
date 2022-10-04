package command

import (
	"fmt"
	"kronos-utils"
	"nym-ws-chat/config"
	"strings"
)

type HelpCmd struct {
	command
}

func NewHelpCmd(name string, argsRequired int) *HelpCmd {
	return &HelpCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *HelpCmd) Execute(cfg *config.Config, args []string) {
	var sb strings.Builder

	sb.WriteString(utils.YELLOW + "USAGE:\n" + utils.RESET)
	sb.WriteString("\t")
	sb.WriteString(utils.RED + args[0] + utils.RESET)
	sb.WriteString(" <" + utils.GREEN + "SUBCOMMAND" + utils.RESET + ">")
	sb.WriteString(" <PARAMS...>\n")
	sb.WriteString(utils.YELLOW + "SUBCOMMANDS:\n" + utils.RESET)
	for _, v := range Values {
		sb.WriteString("\t")
		sb.WriteString(GetHelp(v))
		sb.WriteString("\n")
	}

	fmt.Print(sb.String())

	cmd.command.done = true
}

func (cmd *HelpCmd) GetParams() string {
	return ""
}

func (cmd *HelpCmd) GetDescription() string {
	return "показать меню с подсказками по командам"
}
