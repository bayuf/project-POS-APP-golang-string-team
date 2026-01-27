package adaptor

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"go.uber.org/zap"
)

type Adaptor struct {
	*UserHandler
	*AuthHandler
}

func NewAdaptor(uc *usecase.UseCase, logger *zap.Logger, config *utils.Configuration) *Adaptor {
	return &Adaptor{
		UserHandler: NewUserHandler(uc.UserService, logger, config),
		AuthHandler: NewAuthHandler(uc.AuthService, logger, config),
	}
}
