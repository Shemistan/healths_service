package service

import (
	"github.com/Shemistan/healths_service/internal/models"
)

type IService interface {
	CheckApiUrl()
}

func NewService(cfg *models.Config) IService {
	return &service{
		cfg: cfg,
	}
}

type service struct {
	cfg *models.Config
}
