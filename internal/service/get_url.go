package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Shemistan/healths_service/internal/models"
)

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

func GetApiListWithName() ([]models.ApiWithName, error) {
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
