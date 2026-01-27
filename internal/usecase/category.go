package usecase

import (
	"math"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
)

func (u UseCase) FindAllCategories(page, limit int) ([]entity.MenuCategory, dto.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	categories, total, err := u.repo.CategoryRepo.FindAll(page, limit)
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

	return categories, pagination, nil
}
