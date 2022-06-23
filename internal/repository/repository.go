package repository

import (
	"go.uber.org/zap"
)

type Repository struct {
	Verify
}

type Verify interface {
	VerifyToken(token string, url string) error
}

func NewRepository(logger *zap.SugaredLogger) *Repository {
	return &Repository{
		Verify: newVerifyTokenRepo(logger),
	}
}
