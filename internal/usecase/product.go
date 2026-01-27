package usecase

import (
	"math"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
)

func (u UseCase) FindAllProducts(page, limit int, categoryID *uint) ([]entity.Product, dto.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	products, total, err := u.repo.ProductRepo.FindAll(page, limit, categoryID)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   totalPages,
		TotalRecords: total,
	}

	return products, pagination, nil
}
