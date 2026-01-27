package adaptor

import (
	"net/http"
	"strconv"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	usecase *usecase.UseCase
}

func NewCategoryHandler(usc *usecase.UseCase) *CategoryHandler {
	return &CategoryHandler{
		usecase: usc,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	ctx := c.Request.Context()

	var ctg dto.CreateCategory
	if err := c.ShouldBindJSON(&ctg); err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"invalid request body",
			err.Error(),
		)
		return
	}

	category, err := h.usecase.CreateCategory(ctx, ctg)
	if err != nil {
		status := http.StatusInternalServerError
		message := "failed to create category"

		// validasi name
		if err.Error() == "category name already exists" {
			status = http.StatusConflict
			message = err.Error()
		}

		utils.ResponseFailed(
			c,
			status,
			message,
			err.Error(),
		)
		return
	}

	utils.ResponseSuccess(
		c,
		http.StatusCreated,
		"category created successfully",
		category,
	)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	categories, pagination, err := h.usecase.FindAllCategories(page, limit)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"failed to fetch categories",
			err.Error(),
		)
		return
	}

	utils.ResponsePagination(
		c,
		http.StatusOK,
		"categories retrieved successfully",
		categories,
		pagination,
	)
}

func (h *CategoryHandler) GetCategoryById(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	if idStr == "" {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to fetch category by id",
			"id category is required",
		)
		return
	}

	// 10 = angka biasa, 64 = big size
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to fetch category by id",
			"id category must be a number",
		)
		return
	}

	category, err := h.usecase.GetCategoryByID(ctx, id)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusNotFound,
			"failed to fetch category by id",
			err.Error(),
		)
		return
	}

	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"success fetch category by id",
		category,
	)
}

func (h *CategoryHandler) UpdateCategoryId(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	if idStr == "" {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"id category is required",
			nil,
		)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"id category must be a number",
			err.Error(),
		)
		return
	}

	// bind request body
	var req dto.CreateCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"invalid request body",
			err.Error(),
		)
		return
	}

	updatedCategory, err := h.usecase.UpdateCategoryId(ctx, id, req)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to update category",
			err.Error(),
		)
		return
	}

	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"category updated successfully",
		updatedCategory,
	)
}
