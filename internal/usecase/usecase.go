package usecase

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"go.uber.org/zap"
)

type UseCase struct {
	repo repository.Repository
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger) *UseCase {
	return &UseCase{
		repo: *repo,
	}
}
