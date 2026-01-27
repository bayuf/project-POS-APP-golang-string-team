package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	CategoryRepo CategoryRepository
	ProductRepo  ProductsRepository
}

func NewRepository(db *gorm.DB, logger *zap.Logger) *Repository {
	return &Repository{
		ProductRepo:  NewProductsRepository(db, logger),
		CategoryRepo: NewCategoryRepository(db, logger),
	}
}
