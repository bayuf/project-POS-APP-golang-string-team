package adaptor

import (
	"net/http"
	"strconv"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/usecase"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase *usecase.UseCase
}

func NewProductHandler(usc *usecase.UseCase) *ProductHandler {
	return &ProductHandler{
		usecase: usc,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreateProduct
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"invalid request body",
			err.Error(),
		)
		return
	}

	if err := h.usecase.CreateProduct(ctx, req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "product already exists" {
			status = http.StatusConflict
		}
		utils.ResponseFailed(
			c,
			status,
			err.Error(),
			err.Error(),
		)
		return
	}

	utils.ResponseSuccess(
		c,
		http.StatusCreated,
		"product created successfully",
		nil,
	)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var categoryID *uint
	if c.Query("category_id") != "" {
		id, _ := strconv.Atoi(c.Query("category_id"))
		tmp := uint(id)
		categoryID = &tmp
	}

	products, pagination, err := h.usecase.FindAllProducts(page, limit, categoryID)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"failed to fetch products",
			err.Error(),
		)
		return
	}

	utils.ResponsePagination(
		c,
		http.StatusOK,
		"categories retrieved successfully",
		products,
		pagination,
	)
}

func (h *ProductHandler) GetProductId(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	if idStr == "" {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to fetch product by id",
			"id product is required",
		)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to fetch product by id",
			"id product must be a number",
		)
		return
	}

	product, err := h.usecase.FindProductById(ctx, id)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusNotFound,
			"failed to fetch product by id",
			err.Error(),
		)
		return
	}

	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"success fetch product by id",
		product,
	)
}

func (h *ProductHandler) DeleteProductId(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	if idStr == "" {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to delete product by id",
			"id product is required",
		)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusBadRequest,
			"failed to delete product by id",
			"id product must be a number",
		)
		return
	}

	product, err := h.usecase.FindProductById(ctx, id)
	if err != nil || product == nil {
		utils.ResponseFailed(
			c,
			http.StatusNotFound,
			"failed to delete product by id",
			"product not found",
		)
		return
	}

	// soft delete = is_available
	if err := h.usecase.DeleteProductById(ctx, id); err != nil {
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"failed to delete product by id",
			err.Error(),
		)
		return
	}

	utils.ResponseSuccess(
		c,
		http.StatusOK,
		"success delete product by id",
		product,
	)
}

func (h *ProductHandler) UpdateProductId(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	if idStr == "" {
		utils.ResponseFailed(c, http.StatusBadRequest, "id product is required", nil)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "id product must be a number", nil)
		return
	}

	// request payload
	var req dto.CreateProduct
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		utils.ResponseFailed(c, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	product, err := h.usecase.UpdateProductById(ctx, id, req)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "failed to update product"
		if err.Error() == "product id not found" {
			status = http.StatusNotFound
			msg = err.Error()
		} else if err.Error() == "product name already exists" {
			status = http.StatusConflict
			msg = err.Error()
		}

		utils.ResponseFailed(c, status, msg, err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "product updated successfully", product)
}

func (h *ProductHandler) GetProductsByCategoryName(c *gin.Context) {
	ctx := c.Request.Context()

	name := c.Param("name")
	if name == "" {
		utils.ResponseFailed(c, http.StatusBadRequest, "category name required", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, pagination, err := h.usecase.FindProductsByCategoryName(ctx, name, page, limit)
	if err != nil {
		utils.ResponseFailed(
			c,
			http.StatusInternalServerError,
			"failed to fetch products by category",
			err.Error(),
		)
		return
	}

	utils.ResponsePagination(
		c,
		http.StatusOK,
		"products retrieved successfully",
		products,
		pagination,
	)
}
