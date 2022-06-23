package service

import (
	"client/config"
	"client/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	VerifyToken
}

type VerifyToken interface {
	VerifyToken(token string, url string) error
}

func NewService(repo *repository.Repository, cfg *config.Configs, logger *zap.SugaredLogger) *Service {
	return &Service{
		VerifyToken: NewVerifyTokenService(repo.Verify, cfg, logger),
	}
}
