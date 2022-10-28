package command

import (
	"fmt"
	. "nym-ws-chat/client"
	"nym-ws-chat/client/request"
	"nym-ws-chat/config"
	"strconv"
	"strings"
	"time"
)

type BenchmarkCmd struct {
	command
}

func NewBenchmarkCmd(name string, argsRequired int) *BenchmarkCmd {
	return &BenchmarkCmd{
		command{
			Name:         name,
			ArgsRequired: argsRequired,
		},
	}
}

func (cmd *BenchmarkCmd) Execute(config *config.Config, args []string) {
	client := NewClient(config.Client.Host, config.Client.Port)

	contactNum, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	payloadLength, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println(err)
		return
	}

	benchCount, err := strconv.Atoi(args[4])
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
	text := strings.Repeat("a", payloadLength)

	// Включаем чтение сообщений из сокета
	go client.ReadSocketLoop()

	// Начинаем отсчёт времени
	start := time.Now()

	// Задаём кол-во отправляемых сообщений
	client.Benchmark.N = benchCount
	for i := 0; i < benchCount; i++ {
		writer := client.GetBinaryWriter()
		request.NewSendRequest(writer, false, contact.Address).SetMessage(text).Send()
		writer.Close()
	}

	// Ожидаем получения всех отправленных сообщений
	for client.Benchmark.N > 0 {
	}

	client.Close()

	// Вывод потраченного времени
	fmt.Println("Elapsed time:", time.Since(start))
	fmt.Println("Отправлено", benchCount, "сообщений, размером в", payloadLength, "байт")

	cmd.command.done = true
}

func (cmd *BenchmarkCmd) GetParams() string {
	return "<контакт> <размер> <кол-во>"
}

func (cmd *BenchmarkCmd) GetDescription() string {
	return "отправка сообщений и подсчёт времени работы"
}
