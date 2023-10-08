package service

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Shemistan/healths_service/internal/models"
)

type IService interface {
	CheckApiUrl(ctx context.Context, bot *tgbotapi.BotAPI)
}

func NewService(cfg *models.Config) IService {
	return &service{
		cfg: cfg,
	}
}

type service struct {
	cfg *models.Config
}
