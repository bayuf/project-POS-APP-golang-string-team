package adaptor

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"go.uber.org/zap"
)

type Adaptor struct {
	Category *CategoryHandler
	Products *ProductHandler
}

func NewAdaptor(uc *usecase.UseCase, logger *zap.Logger, config *utils.Configuration) *Adaptor {
	return &Adaptor{
		Category: NewCategoryHandler(uc),
		Products: NewProductHandler(uc),
	}
}
