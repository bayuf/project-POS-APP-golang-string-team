package adaptor

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"go.uber.org/zap"
)

type Adaptor struct {
	CategoryHandler *CategoryHandler
	ProductsHandler *ProductHandler
	*UserHandler
	*AuthHandler
}

func NewAdaptor(uc *usecase.UseCase, logger *zap.Logger, config *utils.Configuration) *Adaptor {
	return &Adaptor{
		CategoryHandler: NewCategoryHandler(uc),
		ProductsHandler: NewProductHandler(uc),
		UserHandler:     NewUserHandler(uc.UserService, logger, config),
		AuthHandler:     NewAuthHandler(uc.AuthService, logger, config),
	}
}
