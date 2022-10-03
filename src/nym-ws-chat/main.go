package main

import (
	"bufio"
	"fmt"
	utils "kronos-utils"
	"nym-ws-chat/client"
	"nym-ws-chat/config"
	"nym-ws-chat/message"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	Config *config.Config
	Client *client.Client
)

func listenConsoleCommands() {
	for !Client.Closed {
		// Набор сообщения
		reader := bufio.NewReader(os.Stdin)
		raw, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			continue
		}

		cmd := strings.Split(string(raw), " ")
		switch cmd[0] {
		case "help":
			fmt.Println(
				"Commands:\n",
				"send <номер_контакта> <сообщение>\n",
				"test <номер_контакта> <сообщение>\n",
				"addcontact <псевдоним> <адрес>\n",
				"delcontact <номер_контакта>\n",
				"list\n",
				"addr",
			)
		case "addr":
			msg := message.NewSelfAddressMessage()
			Client.SendMessage(msg)
		case "list":
			var sb strings.Builder
			sb.WriteString("Список контактов:\n")
			for i, contact := range Config.Contacts {
				sb.WriteString(fmt.Sprintf("#%d %-8s %s\n", i, contact.Alias, contact.Address))
			}
			fmt.Println(sb.String())
		case "addcontact":
			alias := cmd[1]
			address := cmd[2]

			Config.Contacts = append(Config.Contacts, config.Contact{Address: address, Alias: alias})
			Config.Save()
		case "delcontact":
			contactNum, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Println(err)
				continue
			}

			Config.Contacts = append(Config.Contacts[:contactNum], Config.Contacts[contactNum+1:]...)
			Config.Save()
		case "send":
			contactNum, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			text := strings.Join(cmd[2:], " ")

			// Выбор контакта из списка
			if contactNum >= len(Config.Contacts) {
				continue
			}
			contact := Config.Contacts[contactNum]

			// Обправка сообщения
			msg := message.NewOneWayMessage(text, contact.Address, false)
			Client.SendMessage(msg)
		case "test":
			contactNum, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			text := strings.Join(cmd[2:], " ")

			// Выбор контакта из списка
			if contactNum >= len(Config.Contacts) {
				continue
			}
			contact := Config.Contacts[contactNum]

			// Обправка сообщения
			msg := message.NewOneWayMessage(text, contact.Address, true)
			start := time.Now()
			Client.Benchmark.N = 10000
			for i := 0; i < 10000; i++ {
				Client.SendMessage(msg)
			}
			for Client.Benchmark.N != 0 {
			}
			fmt.Println("Time elapsed:", time.Since(start))

		}

		fmt.Println()
	}
}

func run() {
	// MAIN CODE SECTION
	Client = client.NewClient(Config.Client.Host, Config.Client.Port)

	// Мониторинг входящих сообщений (фоновое)
	Client.PrintInputMessages()

	listenConsoleCommands()
}

func init() {
	// Экспорт конфига из встроенного хранилища
	utils.SaveResource("config.yaml", config.BinConfig)

	// Инициализация конфига
	Config = config.NewConfig("config.yaml")
}

func main() {
	interrupt := make(chan os.Signal)      // Channel to listen for interrupt signal to terminate gracefully
	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	go run()

	for {
		select {
		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			fmt.Println("\nTerminate gracefully...")
			Client.Close()
			return
		}
	}
}
