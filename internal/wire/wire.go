package wire

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/adaptor"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Wiring(repo *repository.Repository, logger *zap.Logger, config *utils.Configuration) *gin.Engine {
	// init usecase and adaptor
	uc := usecase.NewUseCase(repo, logger)
	adaptor := adaptor.NewAdaptor(uc, logger, config)

	router := gin.Default() // use default middleware
	r1 := router.Group("/api/v1")
	wireUser(r1, adaptor)

	return router
}

// All Route Here
func wireUser(router *gin.RouterGroup, adaptor *adaptor.Adaptor) {
	router.POST("/users", adaptor.CreateUser)
}
