package main

import (
	"fmt"
	"kronos-utils"
	Command "nym-ws-chat/command"
	"nym-ws-chat/config"
	"os"
	"os/signal"
)

var (
	interrupt chan os.Signal
	done      chan bool

	cfg *config.Config
	cmd Command.Command
)

func run() {
	// MAIN CODE SECTION
	if len(os.Args) < 2 {
		Command.HELP.Execute(cfg, os.Args)
		done <- true
		return
	}

	cmd = Command.GetCommandByName(os.Args[1])
	if cmd == nil {
		fmt.Printf("Команда '%s' не распознана.\n", os.Args[1])
		done <- true
		return
	}

	if !cmd.ValidArgsLength(os.Args) {
		fmt.Printf("Неверное количество аргументов для %s (%d < %d)\n", cmd.GetName(), len(os.Args), cmd.GetRequiredArgsLength())
		done <- true
		return
	}

	cmd.Execute(cfg, os.Args)

	if cmd.IsDone() {
		done <- true
	}
}

func init() {
	// Экспорт конфига из встроенного хранилища
	utils.SaveResource("config.yaml", config.BinConfig)

	// Инициализация конфига
	cfg = config.NewConfig("config.yaml")
}

func main() {
	done = make(chan bool)
	interrupt = make(chan os.Signal)       // Channel to listen for interrupt signal to terminate gracefully
	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	go run()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			fmt.Println("\nTerminate gracefully...")
			cmd.StopExecution()
			return
		}
	}
}
