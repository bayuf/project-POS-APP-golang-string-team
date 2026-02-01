package adaptor

import (
	"net/http"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/middleware"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	usecase *usecase.NotificationService
	logger  *zap.Logger
	config  *utils.Configuration
}

func NewNotificationHandler(uc *usecase.NotificationService, logg *zap.Logger, conf *utils.Configuration) *NotificationHandler {
	return &NotificationHandler{
		usecase: uc,
		logger:  logg,
		config:  conf,
	}
}

func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {
	// pake session auth
	session, _ := middleware.GetAuthUser(c)
	data, err := h.usecase.GetMyNotifications(
		c.Request.Context(),
		session.UserID,
	)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"failed",
			err.Error(),
		)
		return
	}
	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"success",
		data,
	)
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	session, _ := middleware.GetAuthUser(c)

	if err := h.usecase.MarkAsRead(
		c.Request.Context(),
		id,
		session.UserID,
	); err != nil {
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"failed",
			err.Error(),
		)
		return
	}
	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"marked as read",
		nil,
	)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	session, _ := middleware.GetAuthUser(c)

	if err := h.usecase.RemoveNotification(
		c.Request.Context(),
		id,
		session.UserID,
	); err != nil {
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"failed",
			err.Error(),
		)
		return
	}
	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"deleted",
		nil,
	)
}
