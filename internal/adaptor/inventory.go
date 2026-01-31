package adaptor

import (
	"net/http"
	"strconv"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	usecase *usecase.UseCase
}

func NewInventoryHandler(usc *usecase.UseCase) *InventoryHandler {
	return &InventoryHandler{
		usecase: usc,
	}
}

func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	ctx := c.Request.Context()

	var inven dto.CreateInventory
	if err := c.ShouldBindJSON(&inven); err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"invalid request body",
			err.Error(),
		)
		return
	}

	inventory, err := h.usecase.InventoryService.CreateInventory(ctx, inven)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"failed to create inventory",
			err.Error(),
		)
		return
	}

	utils.ResponseSuccess(
		c,
		http.StatusCreated,
		"inventory created successfully",
		inventory,
	)
}

func (h *InventoryHandler) GetAllInventories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	inventories, pagination, err :=
		h.usecase.InventoryService.FindAllInventory(
			c.Request.Context(),
			page,
			limit,
		)

	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"failed to fetch inventories",
			err.Error(),
		)
		return
	}

	utils.ResponsePagination(
		c,
		http.StatusOK,
		"success fetch inventories",
		inventories,
		pagination,
	)
}

func (h *InventoryHandler) FindInventoryDetail(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	if idStr == "" {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to fetch inventories by id",
			"id inventories is required",
		)
		return
	}

	// 10 = angka biasa, 64 = big size
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to fetch inventories by id",
			"id inventories must be a number",
		)
		return
	}

	inven, err := h.usecase.InventoryService.GetDetailInventory(ctx, id)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusNotFound,
			"failed fetched inventory",
			err.Error(),
		)
		return
	}

	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"succesfully get inventory detail",
		inven,
	)
}

func (h *InventoryHandler) UpdateInventoryId(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	if idStr == "" {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"id inventory is required",
			nil,
		)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"id inventory must be a number",
			err.Error(),
		)
		return
	}

	var req dto.UpdateInventory
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"invalid request body",
			err.Error(),
		)
		return
	}

	updateInventory, err := h.usecase.UpdateInventoryId(ctx, id, req)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to update inventory",
			err.Error(),
		)
		return
	}

	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"succesfully update inventory",
		updateInventory,
	)
}
