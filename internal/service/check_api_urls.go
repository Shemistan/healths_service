package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Shemistan/healths_service/internal/models"
	"github.com/Shemistan/healths_service/internal/telegram"
)

func (s *service) CheckApiUrl(ctx context.Context, bot *tgbotapi.BotAPI) {
	apiWithNameList, err := getApiListWithName()
	if err != nil {
		log.Fatal(err)
	}
	// Указываем количество горутин
	goroutineCount := s.cfg.GoroutineCount

	requests := make(chan models.ApiWithName, len(apiWithNameList))
	wg := sync.WaitGroup{}

	go func() {
		for _, api := range apiWithNameList {
			requests <- api
		}
	}()

	defer close(requests)

	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()
			checkUrl(requests, s.cfg.ChatID, bot)
		}()
	}

	go func() {
		wg.Wait()
	}()

	for {
		select {
		case <-time.After(s.cfg.Delay):
			for _, api := range apiWithNameList {
				requests <- api
			}
		case <-time.After(s.cfg.DelayBd):
			log.Println("Add to bd logs")
		case <-ctx.Done():
			return
		}
	}
	//ticker := time.NewTicker(s.cfg.Delay)
	//for range ticker.C {
	//	for _, api := range apiWithNameList {
	//		requests <- api
	//	}
	//}
}

// getApiUrl
// Получаем список урлов из текстового файла
func getApiUrl() ([]string, error) {
	apiList := make([]string, 0)

	file, err := os.Open("config/api_list.txt")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := scanner.Text()
		apiList = append(apiList, url)
	}

	if scanner.Err() != nil {
		fmt.Println("Ошибка чтения файла:", scanner.Err())
	}
	return apiList, nil
}

func getApiListWithName() ([]models.ApiWithName, error) {
	apiList, err := getApiUrl()
	if err != nil {
		return nil, err
	}
	if len(apiList) == 0 {
		return nil, errors.New("empty api list")
	}

	apiWithNameList := make([]models.ApiWithName, 0)

	for _, api := range apiList {
		words := strings.Fields(api)
		apiWithNameList = append(apiWithNameList, models.ApiWithName{
			Name: words[0],
			Url:  words[1],
		})
	}
	return apiWithNameList, nil
}

func checkUrl(requests chan models.ApiWithName, chatID int64, bot *tgbotapi.BotAPI) {
	var status string
	for api := range requests {
		resp, err := http.Get(api.Url)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("[%s] Ошибка мониторинга сервиса %s\n", api.Name, api.Url)
			status = "error"

			msg := fmt.Sprintf("[%s] Ошибка при обращении к сервису: %s", api.Name, api.Url)
			telegram.SendErrorMessage(bot, chatID, msg)
		} else {
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
		fmt.Println(modelForBd)

		if resp != nil {
			resp.Body.Close()
		}
	}
}
