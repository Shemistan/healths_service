package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Shemistan/healths_service/config"
)

func main() {
	log.Println("start monitoring service")
	// init config
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed  to get config:", err.Error())
	}

	// get urls of file
	apiList, err := getApiUrl()
	if err != nil {
		log.Fatal("failed get urls", err)
	}

	// Указываем количество горутин
	goroutineCount := conf.GoroutineCount

	requests := make(chan string, len(apiList))
	//done := make(chan bool)
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
		//done <- true
	}()

	ticker := time.NewTicker(conf.Delay)
	for {
		select {
		//case <-time.After(conf.Delay):
		//	log.Println("wait run after delay monitoring service")
		//time.Sleep(conf.Delay)
		case <-ticker.C:
			for _, api := range apiList {
				requests <- api
			}
			//default:
			//	time.After(conf.Delay)
			//	log.Println("wait run after delay monitoring service")
		}
	}

	//for range ticker.C {
	//	for _, api := range apiList {
	//		resp, err := http.Get(api)
	//		if err != nil || resp.StatusCode != http.StatusOK {
	//			fmt.Printf("Ошибка мониторинга сервиса %s\n", api)
	//		} else {
	//			fmt.Printf("Сервис %s работает исправно\n", api)
	//		}
	//		if resp != nil {
	//			resp.Body.Close()
	//		}
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
