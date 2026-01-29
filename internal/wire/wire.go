package wire

import (
	"net/http"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/adaptor"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/middleware"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Wiring(tx *gorm.DB, repo *repository.Repository, logger *zap.Logger, config *utils.Configuration) *gin.Engine {
	// init usecase and adaptor
	uc := usecase.NewUseCase(repo, logger, tx)
	adaptor := adaptor.NewAdaptor(uc, logger, config)

	// init middleware
	authMW := middleware.NewAuthMiddleware(repo, logger)

	router := gin.Default() // use default middleware

	router.GET("/", func(c *gin.Context) {
		utils.ResponseSuccess(c, http.StatusOK, "OK", gin.H{
			"message": "System POST API Successfully Running",
		})
	})

	r1 := router.Group("/api/v1")

	// Wiring Routes
	wireUser(r1, adaptor, authMW)
	wireAuth(r1, adaptor, authMW)
	wireMenuManagement(r1, adaptor)

	return router
}

// All Route Here
func wireUser(router *gin.RouterGroup, adaptor *adaptor.Adaptor, mw *middleware.AuthMiddleware) {
	users := router.Group("/users")
	users.Use(mw.SessionAuthMiddleware())
	users.GET("/profile", adaptor.GetMyProfile)
	users.PATCH("/profile", adaptor.UpdateMyProfile)
	users.GET("/admins", adaptor.GetAllAdmins)
	users.PATCH("/permissions/:id", mw.RequireRoles("superadmin"), adaptor.UpdateUserPermissions)
	users.Use(mw.RequireRoles("superadmin", "admin"))
	users.GET("", adaptor.GetAllUsers)
	users.GET("/:id", adaptor.GetUserByID)
	users.PUT("/:id", adaptor.UpdateUser)
	users.DELETE("/:id", adaptor.DeleteUser)
	users.POST("", adaptor.CreateUser)
}

func wireAuth(router *gin.RouterGroup, adaptor *adaptor.Adaptor, mw *middleware.AuthMiddleware) {
	auth := router.Group("/auth")
	auth.POST("/login", adaptor.Login)
	auth.POST("/forget-password", adaptor.GetOtpResetPassword)
	auth.POST("/verify-otp", adaptor.GetSessionResetPassword)
	auth.POST("/reset-password", adaptor.ResetPassword)
	auth.Use(mw.SessionAuthMiddleware())
	auth.POST("/logout", adaptor.Logout)
}

func wireMenuManagement(router *gin.RouterGroup, adaptor *adaptor.Adaptor) {
	menu := router.Group("/menu")
	{
		categories := menu.Group("/categories")
		{
			categories.GET("", adaptor.CategoryHandler.GetAllCategories) // path "/"
			categories.GET("/:id", adaptor.CategoryHandler.GetCategoryById)
			categories.PUT("/:id", adaptor.CategoryHandler.UpdateCategoryId)
			categories.DELETE("/:id", adaptor.CategoryHandler.DeleteCategoryById)
			categories.POST("", adaptor.CategoryHandler.CreateCategory)
		}

		products := menu.Group("/products")
		{
			products.GET("", adaptor.ProductsHandler.GetAllProducts)
			products.GET("/:id", adaptor.ProductsHandler.GetProductId)
			products.POST("", adaptor.ProductsHandler.CreateProduct)
			products.PUT("/:id", adaptor.ProductsHandler.UpdateProductId)
			products.DELETE("/:id", adaptor.ProductsHandler.DeleteProductId)
			products.GET("/category/:name", adaptor.ProductsHandler.GetProductsByCategoryName)
		}
	}
}
