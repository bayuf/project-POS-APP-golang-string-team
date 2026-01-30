package usecase

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UseCase struct {
	log  *zap.Logger
	repo repository.Repository

	*UserService
	*AuthService
	*ProductService
	*InventoryService
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger, tx *gorm.DB) *UseCase {
	return &UseCase{
		log:  logger,
		repo: *repo,

		UserService:      NewUserService(repo.UserRepository, logger),
		AuthService:      NewAuthService(repo.AuthRepository, logger, tx),
		ProductService:   NewProductService(repo.ProductRepo, logger, tx),
		InventoryService: NewInventoryService(repo.InventoryRepo, logger, tx),
	}
}
