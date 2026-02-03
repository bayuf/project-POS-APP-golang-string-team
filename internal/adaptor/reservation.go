package adaptor

import (
	"net/http"
	"strconv"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReservationHandler struct {
	usecase *usecase.ReservationService
	logger  *zap.Logger
	config  *utils.Configuration
}

func NewReservationHandler(uc *usecase.ReservationService, logg *zap.Logger, conf *utils.Configuration) *ReservationHandler {
	return &ReservationHandler{
		usecase: uc,
		logger:  logg,
		config:  conf,
	}
}

func (h *ReservationHandler) FindAllReservations(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reservations, pagination, err := h.usecase.FindAllReservation(ctx, page, limit)
	if err != nil {
		h.logger.Error("failed to fetch reservations", zap.Error(err))
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"failed to fetch reservations",
			err.Error(),
		)
		return
	}

	utils.ResponsePagination(
		c,
		http.StatusOK,
		"reservations retrieved successfully",
		reservations,
		pagination,
	)
}

func (h *ReservationHandler) Create(c *gin.Context) {
	var req dto.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"invalid request body",
			err.Error(),
		)
		return
	}

	res, err := h.usecase.CreateReservation(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to create reservation", zap.Error(err))
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to create reservation",
			err.Error(),
		)
		return
	}

	h.logger.Info("reservation created successfully", zap.String("id", res.ID.String()))
	utils.ResponseSuccess(
		c,
		http.StatusCreated,
		"reservation created successfully",
		res,
	)
}

func (h *ReservationHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.usecase.GetReservationDetail(c.Request.Context(), id)
	if err != nil {
		h.logger.Warn("reservation not found", zap.String("id", id))
		utils.ResponseFailed(c, http.StatusNotFound, "reservation not found", nil)
		return
	}

	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"reservation retrieved successfully",
		res,
	)
}

func (h *ReservationHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"invalid request body",
			err.Error(),
		)
		return
	}

	res, err := h.usecase.UpdateReservation(c.Request.Context(), id, req)
	if err != nil {
		h.logger.Error("failed to update reservation", zap.Error(err))
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"update failed",
			err.Error(),
		)
		return
	}

	h.logger.Info("reservation updated", zap.String("id", id))
	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"reservation updated successfully",
		res,
	)
}
