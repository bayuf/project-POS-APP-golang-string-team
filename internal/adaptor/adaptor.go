package adaptor

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"go.uber.org/zap"
)

type Adaptor struct{}

func NewAdaptor(uc *usecase.UseCase, logger *zap.Logger, config *utils.Configuration) *Adaptor {
	return &Adaptor{}
}
