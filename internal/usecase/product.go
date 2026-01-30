package usecase

import (
	"context"
	"errors"
	"math"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductService struct {
	repo   repository.ProductsRepository
	logger *zap.Logger
	tx     *gorm.DB
}

func NewProductService(repo repository.ProductsRepository, logger *zap.Logger, tx *gorm.DB) *ProductService {
	return &ProductService{
		repo:   repo,
		logger: logger,
		tx:     tx,
	}
}

func (u *ProductService) CreateProduct(ctx context.Context, req dto.CreateProduct) error {
	return u.tx.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		repo := repository.NewProductsRepository(tx, u.logger)

		exist, err := repo.IsProductExists(ctx, req.Name)
		if err != nil {
			return err
		}
		if exist != nil {
			return errors.New("product already exists")
		}

		product := &entity.Product{
			Name:        req.Name,
			Price:       req.Price,
			CategoryID:  req.CategoryID,
			IsAvailable: true,
		}

		return repo.Create(ctx, product)
	})
}

func (u UseCase) FindAllProducts(page, limit int, categoryID *uint) ([]entity.Product, dto.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	products, total, err := u.repo.ProductRepo.FindAll(context.Background(), page, limit, categoryID)
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

func (u UseCase) DeleteProductById(ctx context.Context, id int64) error {
	return u.repo.ProductRepo.SoftDeleteById(ctx, id)
}

func (u UseCase) UpdateProductById(ctx context.Context, id int64, req dto.CreateProduct) (*entity.Product, error) {
	existProd, err := u.repo.ProductRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if existProd == nil {
		return nil, errors.New("product id not found")
	}

	checkUniqNameProd, err := u.repo.ProductRepo.IsProductExists(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if checkUniqNameProd != nil && checkUniqNameProd.ID != id {
		return nil, errors.New("product name already exists")
	}

	updateProd := &entity.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Price:       req.Price,
		IsAvailable: true,
	}

	return u.repo.ProductRepo.UpdateProductId(ctx, id, updateProd)
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
