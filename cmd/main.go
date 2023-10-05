package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Shemistan/healths_service/config"
)

func main() {
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

	ticker := time.NewTicker(conf.Delay)

	for range ticker.C {
		for _, api := range apiList {
			resp, err := http.Get(api)
			if err != nil || resp.StatusCode != http.StatusOK {
				fmt.Printf("Ошибка мониторинга сервиса %s\n", api)
			} else {
				fmt.Printf("Сервис %s работает исправно\n", api)
			}
			if resp != nil {
				resp.Body.Close()
			}
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
		return nil, nil
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
