package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/Shemistan/healths_service/internal/models"
)

const LocalRun = true

func NewConfig() (models.Config, error) {
	cfg := models.Config{}
	var err error

	if LocalRun {
		err = godotenv.Load("./dev/local.env")
		if err != nil {
			return models.Config{}, err
		}

		var cfg models.Config
		err = envconfig.Process("", &cfg)
		if err != nil {
			return models.Config{}, err
		}
		return cfg, nil
	}

	envconfig.MustProcess("", &cfg)
	return cfg, nil

}
