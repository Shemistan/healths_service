package models

import "time"

type Config struct {
	App            APP           `envconfig:"APP"`
	Delay          time.Duration `envconfig:"DELAY"`
	GoroutineCount int           `envconfig:"GOROUTINE_COUNT"`
}

type APP struct {
	Port string `envconfig:"PORT"`
}
