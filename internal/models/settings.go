package models

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Settings struct {
	ChatID int64
	Urls   []ApiWithName
	Bot    *tgbotapi.BotAPI
}
