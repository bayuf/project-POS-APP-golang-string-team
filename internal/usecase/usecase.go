package usecase

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
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
	*OrderService
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger, tx *gorm.DB, config *utils.Configuration, emailJob chan<- utils.EmailJob) *UseCase {
	return &UseCase{
		log:  logger,
		repo: *repo,

		ProductService:   NewProductService(repo.ProductRepo, logger, tx),
		InventoryService: NewInventoryService(repo.InventoryRepo, logger, tx),
		UserService:      NewUserService(repo.UserRepository, logger, emailJob),
		AuthService:      NewAuthService(repo.AuthRepository, logger, tx, emailJob),
		OrderService:     NewOrderService(repo.OrderRepository, logger, tx),
	}
}
