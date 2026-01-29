package adaptor

import (
	"net/http"
	"strconv"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/middleware"
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

func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	user := dto.UpdateUser{}
	if err := c.BindJSON(&user); err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "failed", err.Error())
		return
	}

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

	if err := h.uc.UpdateUserData(ctx, id, user); err != nil {
		utils.ResponseFailed(c, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "success", nil)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
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

	if err := h.uc.DeleteUserByID(ctx, id); err != nil {
		utils.ResponseFailed(c, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "success", nil)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	pageStr := c.Query("page")
	sortBy := c.Query("sortby")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "failed", err.Error())
		return
	}

	filter := dto.UserFilterRequest{
		Page:   page,
		Limit:  h.config.Limit,
		SortBy: sortBy,
	}

	users, pagination, err := h.uc.GetAllUser(ctx, filter)
	if err != nil {
		utils.ResponseFailed(c, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponsePagination(c, http.StatusOK, "success", users, *pagination)
}

func (h *UserHandler) GetMyProfile(c *gin.Context) {
	ctx := c.Request.Context()

	// Must Login
	user, ok := middleware.GetAuthUser(c)
	if !ok {
		utils.ResponseFailed(c, http.StatusUnauthorized, "failed", "user not found")
		return
	}

	userProfile, err := h.uc.GetUserProfile(ctx, user.UserID)
	if err != nil {
		utils.ResponseFailed(c, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "success", userProfile)
}
