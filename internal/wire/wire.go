package wire

import (
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
	r1 := router.Group("/api/v1")

	// Wiring Routes
	wireUser(r1, adaptor, authMW)
	wireAuth(r1, adaptor, authMW)

	return router
}

// All Route Here
func wireUser(router *gin.RouterGroup, adaptor *adaptor.Adaptor, mw *middleware.AuthMiddleware) {
	users := router.Group("/users")
	users.Use(mw.SessionAuthMiddleware(), mw.RequireRoles("superadmin", "admin"))
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
