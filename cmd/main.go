package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Shemistan/healths_service/config"
	"github.com/Shemistan/healths_service/internal/service"
)

func main() {
	log.Println("start monitoring service")
	// init config
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed  to get config:", err.Error())
	}

	// Создание бота
	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		termCh := make(chan os.Signal, 1)
		signal.Notify(termCh, os.Interrupt, syscall.SIGINT)
		<-termCh
		log.Println("Shutdown...")
		cancel()
	}()

	serv := service.NewService(&conf)

	apiWithNameList, err := service.GetApiListWithName()
	if err != nil {
		log.Fatal(err)
	}

	serv.CheckApiUrl(ctx, bot, apiWithNameList)
}
