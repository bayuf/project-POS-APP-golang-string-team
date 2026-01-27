package repository

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductsRepository interface {
	Create(prod *entity.Product) error
	FindAll(page, limit int, categoryID *uint) ([]entity.Product, int64, error)
}

type productsRepo struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewProductsRepository(db *gorm.DB, log *zap.Logger) ProductsRepository {
	return &productsRepo{
		DB:  db,
		Log: log,
	}
}

func (r *productsRepo) Create(prod *entity.Product) error {
	return r.DB.Create(prod).Error
}

func (r *productsRepo) FindAll(page, limit int, categoryID *uint) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64 // record pagination

	offset := (page - 1) * limit

	query := r.DB.Model(&entity.Product{})

	// filter by category if provided
	if categoryID != nil && *categoryID > 0 {
		query = query.Where("category_id = ?", *categoryID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Category").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
