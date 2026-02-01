package adaptor

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"go.uber.org/zap"
)

type Adaptor struct {
	CategoryHandler *CategoryHandler
	ProductsHandler *ProductHandler
	*NotificationHandler
	*InventoryHandler
	*ReservationHandler
	*UserHandler
	*AuthHandler
}

func NewAdaptor(uc *usecase.UseCase, logger *zap.Logger, config *utils.Configuration) *Adaptor {
	return &Adaptor{
		CategoryHandler:     NewCategoryHandler(uc),
		ProductsHandler:     NewProductHandler(uc),
		InventoryHandler:    NewInventoryHandler(uc),
		NotificationHandler: NewNotificationHandler(uc.NotificationService, logger, config),
		ReservationHandler:  NewReservationHandler(uc.ReservationService, logger, config),
		UserHandler:         NewUserHandler(uc.UserService, logger, config),
		AuthHandler:         NewAuthHandler(uc.AuthService, logger, config),
	}
}
