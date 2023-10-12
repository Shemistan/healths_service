package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Shemistan/healths_service/internal/models"
	"github.com/Shemistan/healths_service/internal/telegram"
)

func (s *service) CheckApiUrl(ctx context.Context, bot *tgbotapi.BotAPI) {
	apiWithNameList, err := GetApiListWithName()
	if err != nil {
		log.Fatal(err)
	}
	// Указываем количество горутин
	//goroutineCount := s.cfg.GoroutineCount
	//requests := make(chan models.ApiWithName, len(apiWithNameList))

	stg := models.Settings{
		ChatID: s.cfg.ChatID,
		Urls:   apiWithNameList,
		Bot:    bot,
	}

	for {
		select {
		case <-time.After(s.cfg.Delay):
			checkServices(stg)
		case <-time.After(s.cfg.DelayBd):
			log.Println("Add to bd logs")
		case <-ctx.Done():
			return
		}
	}
}

func checkServices(st models.Settings) {
	var wg sync.WaitGroup
	results := make(chan models.ApiWithName)

	for _, api := range st.Urls {
		wg.Add(1)
		go checkService(api, &wg, st, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	for api := range results {
		log.Println(api)
	}
}

func checkService(api models.ApiWithName, wg *sync.WaitGroup, st models.Settings, results chan<- models.ApiWithName) {
	var status string
	defer wg.Done()

	resp, err := http.Get(api.Url)
	if err != nil || resp.StatusCode != http.StatusOK {
		results <- api
		log.Printf("[%s] Ошибка мониторинга сервиса %s\n", api.Name, api.Url)
		status = "error"

		msg := fmt.Sprintf("[%s] Ошибка при обращении к сервису: %s", api.Name, api.Url)
		telegram.SendErrorMessage(st.Bot, st.ChatID, msg)
	} else {
		results <- api
		log.Printf("Сервис %s работает исправно\n", api)
		status = "success"
	}
	///TODO add to chan for bd
	modelForBd := models.ServiceCheck{
		Name:     api.Name,
		Url:      api.Url,
		CreateAt: time.Now(),
		Status:   status,
	}
	_ = modelForBd
	//fmt.Println(modelForBd)

	if resp != nil {
		resp.Body.Close()
	}
}
