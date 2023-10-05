package models

import "time"

type Config struct {
	App   APP           `envconfig:"APP"`
	Delay time.Duration `envconfig:"DELAY"`
}

type APP struct {
	Port string `envconfig:"PORT"`
}
