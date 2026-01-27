package wire

import (
	"net/http"

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

	router.GET("/", func(c *gin.Context) {
		utils.ResponseSuccess(c, http.StatusOK, "OK", gin.H{
			"message": "System POST API Successfully Running",
		})
	})

	r1 := router.Group("/api/v1")
	wireUser(r1, adaptor)
	wireMenuManagement(r1, adaptor)

	return router
}

// All Route Here
func wireUser(router *gin.RouterGroup, adaptor *adaptor.Adaptor) {
	router.GET("/users/", adaptor.GetAllUsers)
	router.GET("/users/:id", adaptor.GetUserByID)
	router.PUT("/users/:id", adaptor.UpdateUser)
	router.DELETE("/users/:id", adaptor.DeleteUser)
	router.POST("/users", adaptor.CreateUser)
}

func wireMenuManagement(router *gin.RouterGroup, adaptor *adaptor.Adaptor) {
	menu := router.Group("/menu")
	{
		categories := menu.Group("/categories")
		{
			categories.GET("", adaptor.Category.GetAllCategories) // path "/"
		}

		products := menu.Group("/products")
		{
			products.GET("", adaptor.Products.GetAllProducts)
		}
	}
}
