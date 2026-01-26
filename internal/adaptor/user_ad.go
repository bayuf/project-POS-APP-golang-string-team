package adaptor

import (
	"net/http"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	uc     *usecase.UserService
	logger *zap.Logger
	config *utils.Configuration
}

func NewUserHandler(uc *usecase.UserService, log *zap.Logger, conf *utils.Configuration) *UserHandler {
	return &UserHandler{
		uc:     uc,
		logger: log,
		config: conf,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	user := dto.CreateUser{}
	if err := c.BindJSON(&user); err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "failed", err.Error())
		return
	}

	if err := h.uc.CreateUser(ctx, user); err != nil {
		utils.ResponseFailed(c, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "success", nil)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	if idStr == "" {
		utils.ResponseFailed(c, http.StatusBadRequest, "failed", "id is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "failed", err.Error())
		return
	}

	user, err := h.uc.GetUserByID(ctx, id)
	if err != nil {
		utils.ResponseFailed(c, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "success", user)
}
