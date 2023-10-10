package internal

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"strconv"
	"strings"
	"time"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig
	data         map[int64]*Runner
	fileName     string
}

func (b *Bot) Init(fileName string) {
	b.fileName = fileName
	buff, _ := os.ReadFile(fileName)
	rawData := strings.Split(string(buff), "\n")
	data := make(map[int64][]string)

	b.bot, _ = tgbotapi.NewBotAPI(rawData[0])
	b.data = make(map[int64]*Runner)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	if len(rawData) > 1 {
		for i := 1; i < len(rawData); i++ {
			arr := strings.Split(rawData[i], " ")
			chatID, _ := strconv.ParseInt(arr[1], 10, 64)
			data[chatID] = append(data[chatID], arr[2])
		}

		for chatID := range data {
			b.data[chatID] = NewRunner(false)
			b.data[chatID].Add(data[chatID])
		}
	}
}

func (b *Bot) Send(msg string, chatID int64) error {
	_, err := b.bot.Send(tgbotapi.NewMessage(chatID, msg))
	return err
}

func (b *Bot) Receive() {
	updates, _ := b.bot.GetUpdatesChan(b.updateConfig)
	result := make(chan string)
	defer close(result)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "add":
				if b.data[update.Message.Chat.ID] == nil {
					b.data[update.Message.Chat.ID] = NewRunner(false)
				}
				b.data[update.Message.Chat.ID].Add([]string{strings.ReplaceAll(update.Message.Text, "/add ", "")})
			case "remove":
				b.data[update.Message.Chat.ID].Remove(strings.ReplaceAll(update.Message.Text, "/remove ", ""))
			case "run":
				{
					timeout, _ := strconv.Atoi(strings.ReplaceAll(update.Message.Text, "/run ", ""))
					go b.data[update.Message.Chat.ID].Start(result, time.Duration(timeout))
					reading := func() {
						for {
							if err := b.Send(<-result, update.Message.Chat.ID); err != nil {
								fmt.Println(err)
							}
						}
					}
					go reading()
				}
			case "stop":
				b.data[update.Message.Chat.ID].Stop()
			}
		} else {
			if err := b.Send("Use commands to operate with the bot", update.Message.Chat.ID); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (b *Bot) Stop() {
	b.bot.StopReceivingUpdates()
	for chatID := range b.data {
		b.data[chatID].Stop()
		for _, url := range b.data[chatID].Export() {
			if err := os.WriteFile(b.fileName, []byte(strconv.FormatInt(chatID, 10)+" "+url+"\n"), 0644); err != nil {
				fmt.Println(err)
			}
		}
	}
}
