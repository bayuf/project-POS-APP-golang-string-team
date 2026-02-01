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
	*ReservationService
	*NotificationService
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger, tx *gorm.DB, config *utils.Configuration, emailJob chan<- utils.EmailJob) *UseCase {
	return &UseCase{
		log:  logger,
		repo: *repo,

		UserService:         NewUserService(repo.UserRepository, logger, emailJob),
		AuthService:         NewAuthService(repo.AuthRepository, logger, tx, emailJob),
		ProductService:      NewProductService(repo.ProductRepo, logger, tx),
		InventoryService:    NewInventoryService(repo.InventoryRepo, logger, tx),
		ReservationService:  NewReservationService(repo.ReservationRepo, logger, tx),
		NotificationService: NewNotificationService(repo.NotifRepo, logger),
	}
}
