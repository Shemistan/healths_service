package main

import (
	"log"

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

	serv := service.NewService(&conf)
	serv.CheckApiUrl()
}
