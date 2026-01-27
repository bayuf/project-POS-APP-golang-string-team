package adaptor

import (
	"net/http"
	"strconv"

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
