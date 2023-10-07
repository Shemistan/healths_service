package service

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func (s *service) CheckApiUrl() {
	// get urls of file
	apiList, err := getApiUrl()
	if err != nil {
		log.Fatal("failed get urls", err)
	}

	// Указываем количество горутин
	goroutineCount := s.cfg.GoroutineCount

	requests := make(chan string, len(apiList))
	wg := sync.WaitGroup{}

	go func() {
		for _, api := range apiList {
			requests <- api
		}
	}()

	defer close(requests)

	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done()
			for api := range requests {
				resp, err := http.Get(api)
				if err != nil || resp.StatusCode != http.StatusOK {
					log.Printf("Ошибка мониторинга сервиса %s\n", api)
				} else {
					log.Printf("Сервис %s работает исправно\n", api)
				}

				if resp != nil {
					resp.Body.Close()
				}
			}
		}()
	}

	go func() {
		wg.Wait()
	}()

	ticker := time.NewTicker(s.cfg.Delay)
	//for {
	//	select {
	//	case <-ticker.C:
	//		for _, api := range apiList {
	//			requests <- api
	//		}
	//	}
	//}
	for range ticker.C {
		for _, api := range apiList {
			requests <- api
		}
	}
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
