package adaptor

import (
	"net/http"
	"strconv"

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
