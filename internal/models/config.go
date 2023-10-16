package models

import "time"

type Config struct {
	App             APP           `envconfig:"APP"`
	WorkerTimeOut   time.Duration `envconfig:"WORKER_TIME_OUT"`
	WorkerTimeOutBd time.Duration `envconfig:"WORKER_TIME_OUT_BD"`
	GoroutineCount  int           `envconfig:"GOROUTINE_COUNT"`
	Token           string        `envconfig:"TOKEN"`
	ChatID          int64         `envconfig:"CHAT_ID"`
}

type APP struct {
	Port string `envconfig:"PORT"`
}
