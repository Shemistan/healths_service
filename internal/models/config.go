package models

import "time"

type Config struct {
	App            APP           `envconfig:"APP"`
	Delay          time.Duration `envconfig:"DELAY"`
	DelayBd        time.Duration `envconfig:"DELAY_BD"`
	GoroutineCount int           `envconfig:"GOROUTINE_COUNT"`
	Token          string        `envconfig:"TOKEN"`
	ChatID         int64         `envconfig:"CHAT_ID"`
}

type APP struct {
	Port string `envconfig:"PORT"`
}
