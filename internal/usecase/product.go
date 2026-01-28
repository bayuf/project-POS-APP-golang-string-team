package usecase

import (
	"context"
	"errors"
	"math"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
)

func (u UseCase) CreateProduct(ctx context.Context, req dto.CreateProduct) (*entity.Product, error) {

	category, err := u.repo.CategoryRepo.FindById(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	product := &entity.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Price:       req.Price,
		IsAvailable: true,
	}

	if req.IsAvailable != nil {
		product.IsAvailable = *req.IsAvailable
	}

	if err := u.repo.ProductRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

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

func (u UseCase) FindProductById(ctx context.Context, id int64) (*dto.Product, error) {
	product, err := u.repo.ProductRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("product id not found")
	}

	var category *dto.Category
	if product.Category != nil {
		category = &dto.Category{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		}
	}

	return &dto.Product{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		IsAvailable: product.IsAvailable,
		Category:    category,
	}, nil
}

func (u UseCase) FindProductsByCategoryName(ctx context.Context, name string, page, limit int) ([]entity.Product, dto.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	products, total, err := u.repo.ProductRepo.FindByCategoryName(
		ctx,
		name,
		page,
		limit,
	)
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
