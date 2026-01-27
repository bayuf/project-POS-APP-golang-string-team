package usecase

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"go.uber.org/zap"
)

type UseCase struct {
	*UserService
	*AuthService
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger) *UseCase {
	return &UseCase{
		UserService: NewUserService(repo.UserRepository, logger),
		AuthService: NewAuthService(repo.AuthRepository, logger),
	}
}
