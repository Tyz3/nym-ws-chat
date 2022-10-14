package command

import (
	"fmt"
	utils "kronos-utils"
	"nym-ws-chat/client"
	"nym-ws-chat/config"
)

type command struct {
	Name         string
	ArgsRequired int
	done         bool
	client       *client.Client
}

type Command interface {
	ValidArgsLength(args []string) bool
	Execute(config *config.Config, args []string)
	GetName() string
	GetParams() string
	GetDescription() string
	IsDone() bool
	StopExecution()
	GetRequiredArgsLength() int
}

var (
	HELP       Command = NewHelpCmd("help", 2)
	ADDCONTACT Command = NewAddContactCmd("addcontact", 4)
	DELCONTACT Command = NewDelContactCmd("delcontact", 3)
	SEND       Command = NewSendCmd("send", 4)
	SENDF      Command = NewSendFCmd("sendf", 4)
	LISTEN     Command = NewListenCmd("listen", 2)
	LIST       Command = NewListCmd("list", 2)
	ADDR       Command = NewAddrCmd("addr", 2)
	BENCHMARK  Command = NewBenchmarkCmd("benchmark", 5)
	REPLY      Command = NewReplyCmd("reply", 4)
	REPLYF     Command = NewReplyFCmd("replyf", 4)

	Values = []Command{HELP, ADDCONTACT, DELCONTACT, SEND, SENDF, LISTEN, LIST, ADDR, BENCHMARK, REPLY, REPLYF}
)

// ValidArgsLength Абстрактный метод
func (c *command) ValidArgsLength(args []string) bool {
	return c.ArgsRequired <= len(args)
}

// GetName Абстрактный метод
func (c *command) GetName() string {
	return c.Name
}

// IsDone Абстрактный метод
func (c *command) IsDone() bool {
	return c.done
}

// StopExecution Абстрактный метод
func (c *command) StopExecution() {
	if c.client != nil && !c.client.Closed {
		c.client.Close()
	}
}

// GetHelp Статичный метод
func GetHelp(cmd Command) string {
	return fmt.Sprintf(
		"%s%-12s%s %-32s %s",
		utils.GREEN, cmd.GetName(), utils.RESET,
		cmd.GetParams(),
		cmd.GetDescription(),
	)
}

// GetRequiredArgsLength Абстрактный метод
func (c *command) GetRequiredArgsLength() int {
	return c.ArgsRequired
}

func GetCommandByName(cmdName string) Command {
	for _, v := range Values {
		if cmdName == v.GetName() {
			return v
		}
	}
	return nil
}
