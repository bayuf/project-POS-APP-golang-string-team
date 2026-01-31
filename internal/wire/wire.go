package wire

import (
	"net/http"
	"sync"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/adaptor"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/middleware"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Route *gin.Engine
	Stop  chan struct{}
	WG    *sync.WaitGroup
}

func Wiring(tx *gorm.DB, repo *repository.Repository, logger *zap.Logger, config *utils.Configuration) *App {
	// init worker
	emailJobs := make(chan utils.EmailJob, 10) // BUFFER
	stop := make(chan struct{})
	wg := &sync.WaitGroup{}
	emailUC := usecase.NewEmailService(logger, config)

	utils.StartEmailWorkers(4, emailJobs, stop, wg, emailUC)

	// init usecase and adaptor
	uc := usecase.NewUseCase(repo, logger, tx, config, emailJobs)
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
	wireInventory(r1, adaptor)

	return &App{
		Route: router,
		Stop:  stop,
		WG:    wg,
	}
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
	auth.POST("/reset-password", adaptor.GetOtpResetPassword)
	auth.POST("/verify-otp", adaptor.GetSessionResetPassword)
	auth.POST("/update-password", adaptor.ResetPassword)
	auth.Use(mw.SessionAuthMiddleware())
	auth.POST("/logout", adaptor.Logout)
}

// belum ada middleware
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

// belum ada middleware
func wireInventory(router *gin.RouterGroup, adaptor *adaptor.Adaptor) {
	inven := router.Group("/inventories")

	inven.POST("", adaptor.InventoryHandler.CreateInventory)
	inven.GET("", adaptor.InventoryHandler.GetAllInventories)

	// /inventories/search?product_name=garlic&category_id=1&min_stock=10&max_stock=10
	inven.GET("/search", adaptor.InventoryHandler.SearchInventories)

	inven.GET("/:id", adaptor.InventoryHandler.FindInventoryDetail)
	inven.PUT("/:id", adaptor.InventoryHandler.UpdateInventoryId)
	inven.DELETE("/:id", adaptor.InventoryHandler.DeleteInventoryById)
}
