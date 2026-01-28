package usecase

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UseCase struct {
	repo repository.Repository
	*UserService
	*AuthService
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger, tx *gorm.DB) *UseCase {
	return &UseCase{
		repo: *repo,
		UserService: NewUserService(repo.UserRepository, logger),
		AuthService: NewAuthService(repo.AuthRepository, logger, tx),
	}
}
