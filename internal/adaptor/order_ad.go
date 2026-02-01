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

type OrderHandler struct {
	uc     *usecase.OrderService
	logger *zap.Logger
	config *utils.Configuration
}

func NewOrderHandler(uc *usecase.OrderService, logger *zap.Logger, config *utils.Configuration) *OrderHandler {
	return &OrderHandler{
		uc:     uc,
		logger: logger,
		config: config,
	}
}

func (h *OrderHandler) Order(c *gin.Context) {
	ctx := c.Request.Context()
	order := dto.Order{}
	if err := c.ShouldBindJSON(&order); err != nil {
		h.logger.Error("failed to bind json", zap.Error(err))
		utils.ResponseFailed(c, http.StatusBadRequest, "failed", err.Error())
		return
	}

	orderDetail, err := h.uc.CreateOrder(ctx, order)
	if err != nil {
		h.logger.Error("failed to create order", zap.Error(err))
		utils.ResponseFailed(c, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "success", orderDetail)

}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	ctx := c.Request.Context()
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.logger.Error("failed to parse order id", zap.Error(err))
		utils.ResponseFailed(c, http.StatusBadRequest, "failed", err.Error())
		return
	}

	paymentMethod := struct {
		PaymentMethodID int64 `json:"payment_method" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&paymentMethod); err != nil {
		h.logger.Error("failed to bind json", zap.Error(err))
		utils.ResponseFailed(c, http.StatusBadRequest, "failed", err.Error())
		return
	}

	orderDetail, err := h.uc.PayOrder(ctx, paymentMethod.PaymentMethodID, orderID)
	if err != nil {
		h.logger.Error("failed to pay order", zap.Error(err))
		utils.ResponseFailed(c, http.StatusInternalServerError, "failed", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "success", orderDetail)

}
