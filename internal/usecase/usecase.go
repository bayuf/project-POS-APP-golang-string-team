package usecase

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UseCase struct {
	*UserService
	*AuthService
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger, tx *gorm.DB) *UseCase {
	return &UseCase{
		UserService: NewUserService(repo.UserRepository, logger),
		AuthService: NewAuthService(repo.AuthRepository, logger, tx),
	}
}
