package command

import (
	"fmt"
	"nym-ws-chat/client"
	"nym-ws-chat/config"
	"nym-ws-chat/message"
	"strings"
	"time"
)

func Help(exe string) {
	exe = strings.Replace(exe, "./", "  ", 1)
	msg := `Commands:
{exe} send <номер_контакта> <сообщение>
{exe} test <номер_контакта> <размер_сообщения> <количество_сообщений>
{exe} addcontact <имя> <адрес>
{exe} delcontact <номер_контакта>
{exe} list
{exe} addr
{exe} read`

	fmt.Println(strings.ReplaceAll(msg, "{exe}", exe))
}

func Addr(client *client.Client) {
	msg := message.NewSelfAddressMessage()
	client.SendMessage(msg)
}

func List(contacts []config.Contact) {
	var sb strings.Builder
	sb.WriteString("Список контактов:\n")
	for i, contact := range contacts {
		sb.WriteString(fmt.Sprintf("#%d %-8s %s\n", i, contact.Alias, contact.Address))
	}
	fmt.Println(sb.String())
}

func AddContact(cfg *config.Config, alias string, address string) {
	cfg.Contacts = append(cfg.Contacts, config.Contact{Address: address, Alias: alias})
	cfg.Save()
}

func DelContact(cfg *config.Config, contactNum int) {
	cfg.Contacts = append(cfg.Contacts[:contactNum], cfg.Contacts[contactNum+1:]...)
	cfg.Save()
}

func Send(cfg *config.Config, client *client.Client, contactNum int, text string) {
	// Выбор контакта из списка
	if contactNum >= len(cfg.Contacts) {
		fmt.Println("Не найден контакт под номером", contactNum)
		return
	}
	contact := cfg.Contacts[contactNum]

	// Обправка сообщения
	msg := message.NewOneWayMessage(text, contact.Address, false)
	client.SendMessage(msg)
}

func Test(cfg *config.Config, client *client.Client, contactNum int, payloadLength int, benchCount int) {
	ReadMessages(client)
	// Выбор контакта из списка
	if contactNum >= len(cfg.Contacts) {
		fmt.Println("Не найден контакт под номером", contactNum)
		return
	}
	contact := cfg.Contacts[contactNum]

	text := strings.Repeat("a", payloadLength)

	// Обправка сообщения
	msg := message.NewOneWayMessage(text, contact.Address, true)
	start := time.Now()
	client.Benchmark.N = benchCount
	for i := 0; i < benchCount; i++ {
		client.SendMessage(msg)
	}
	for client.Benchmark.N > 0 {
	}
	fmt.Println("Time elapsed:", time.Since(start))
}

func ReadMessages(client *client.Client) {
	channel := make(chan string, 10) // Канал для пересылки сообщений между горутинами
	go client.ReadSocket(channel)
	go client.StartPrint(channel)
}
