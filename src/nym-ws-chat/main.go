package main

import (
	"fmt"
	utils "kronos-utils"
	"nym-ws-chat/client"
	"nym-ws-chat/command"
	"nym-ws-chat/config"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

var (
	interrupt chan os.Signal
	done      chan bool
	Config    *config.Config
	Client    *client.Client
)

func run() {
	// MAIN CODE SECTION
	exe := os.Args[0]
	if len(os.Args) < 2 {
		command.Help(exe)
		done <- true
		return
	}

	cmd := os.Args[1]
	if cmd == "help" {
		command.Help(exe)
		done <- true
	} else if cmd == "addr" {
		Client = client.NewClient(Config.Client.Host, Config.Client.Port)
		command.ReadMessages(Client)
		command.Addr(Client)
	} else if cmd == "list" {
		command.List(Config.Contacts)
		done <- true
	} else if cmd == "addcontact" {
		if len(os.Args) < 4 {
			fmt.Println("Недостаточно аргументов для команды", cmd)
			done <- true
			return
		}

		alias := os.Args[2]
		address := os.Args[3]

		command.AddContact(Config, alias, address)
		command.List(Config.Contacts)
		done <- true
	} else if cmd == "delcontact" {
		if len(os.Args) < 3 {
			fmt.Println("Недостаточно аргументов для команды", cmd)
			done <- true
			return
		}

		contactNum, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		command.DelContact(Config, contactNum)
		command.List(Config.Contacts)
		done <- true
	} else if cmd == "send" {
		if len(os.Args) < 4 {
			fmt.Println("Недостаточно аргументов для команды", cmd)
			done <- true
			return
		}
		Client = client.NewClient(Config.Client.Host, Config.Client.Port)

		contactNum, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		text := strings.Join(os.Args[3:], " ")

		command.Send(Config, Client, contactNum, text)
		Client.Close()
		done <- true
	} else if cmd == "test" {
		// test <номер_контакта> <размер_сообщения>
		if len(os.Args) < 5 {
			fmt.Println("Недостаточно аргументов для команды", cmd)
			done <- true
			return
		}
		Client = client.NewClient(Config.Client.Host, Config.Client.Port)

		contactNum, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		payloadLength, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println(err)
			return
		}

		benchCount, err := strconv.Atoi(os.Args[4])
		if err != nil {
			fmt.Println(err)
			return
		}

		command.Test(Config, Client, contactNum, payloadLength, benchCount)
		Client.Close()
		fmt.Println("Отправлено", benchCount, "сообщений, размером в", payloadLength, "байт")
		done <- true
	} else if cmd == "read" {
		Client = client.NewClient(Config.Client.Host, Config.Client.Port)
		// test <номер_контакта> <размер_сообщения>
		command.ReadMessages(Client)
	} else {
		command.Help(exe)
		done <- true
	}
}

func init() {
	// Экспорт конфига из встроенного хранилища
	utils.SaveResource("config.yaml", config.BinConfig)

	// Инициализация конфига
	Config = config.NewConfig("config.yaml")
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
			if Client != nil {
				Client.Close()
			}
			return
		}
	}
}
