package adaptor

import (
	"net/http"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	uc     *usecase.AuthService
	logger *zap.Logger
}

func NewAuthHandler(uc *usecase.AuthService, logger *zap.Logger, config *utils.Configuration) *AuthHandler {
	return &AuthHandler{
		uc:     uc,
		logger: logger,
	}
}

func (ad *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	loginData := dto.Login{}

	if err := c.BindJSON(&loginData); err != nil {
		ad.logger.Error("failed to bind json", zap.Error(err))
		utils.ResponseFailed(c, http.StatusBadRequest, "failed", err.Error())
		return
	}

	session, err := ad.uc.Login(ctx, loginData)
	if err != nil {
		ad.logger.Error("failed to login", zap.Error(err))
		utils.ResponseFailed(c, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "success", session)
}
