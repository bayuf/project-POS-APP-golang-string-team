package repository

import (
	"context"
	"errors"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductsRepository interface {
	Create(ctx context.Context, prod *entity.Product) error
	FindAll(ctx context.Context, page, limit int, categoryID *uint) ([]entity.Product, int64, error)
	FindById(ctx context.Context, id int64) (*entity.Product, error)
	UpdateProductId(ctx context.Context, id int64, prod *entity.Product) (*entity.Product, error)
	SoftDeleteById(ctx context.Context, id int64) error

	IsProductExists(ctx context.Context, name string) (*entity.Product, error)
	FindByCategoryName(ctx context.Context, name string, page, limit int) ([]entity.Product, int64, error)
	HasActiveProducts(ctx context.Context, categoryID int64) (bool, error)
}

type productsRepo struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewProductsRepository(db *gorm.DB, logger *zap.Logger) ProductsRepository {
	return &productsRepo{
		DB:  db,
		Log: logger,
	}
}

func (r *productsRepo) Create(ctx context.Context, prod *entity.Product) error {
	return r.DB.WithContext(ctx).Create(prod).Error
}

func (r *productsRepo) FindAll(ctx context.Context, page, limit int, categoryID *uint) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64 // record pagination

	offset := (page - 1) * limit

	// ambil yang true, jadi yang false tdk masuk ke response
	query := r.DB.WithContext(ctx).Model(&entity.Product{}).Where("is_available = ?", true)

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

func (r *productsRepo) FindById(ctx context.Context, id int64) (*entity.Product, error) {
	var prod entity.Product

	err := r.DB.WithContext(ctx).
		Preload("Category").
		Preload("Inventory").
		First(&prod, "id = ? AND is_available = ?", id, true).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &prod, nil
}

func (r *productsRepo) UpdateProductId(ctx context.Context, id int64, prod *entity.Product) (*entity.Product, error) {
	updates := map[string]interface{}{
		"name":         prod.Name,
		"price":        prod.Price,
		"category_id":  prod.CategoryID,
		"is_available": prod.IsAvailable,
	}

	res := r.DB.WithContext(ctx).
		Model(&entity.Product{}).
		Where("id = ?", id).
		Updates(updates)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return r.FindById(ctx, id)
}

func (r *productsRepo) SoftDeleteById(ctx context.Context, id int64) error {
	now := time.Now()

	res := r.DB.WithContext(ctx).
		Model(&entity.Product{}).
		Where("id = ? AND is_available = ?", id, true).
		Updates(map[string]interface{}{
			"is_available": false,
			"deleted_at":   &now,
		})

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

// filter category by name
func (r *productsRepo) FindByCategoryName(ctx context.Context, name string, page, limit int) ([]entity.Product, int64, error) {

	var products []entity.Product
	var total int64

	offset := (page - 1) * limit

	query := r.DB.WithContext(ctx).
		Model(&entity.Product{}).
		Joins("JOIN menu_categories mc ON mc.id = products.category_id").
		Where("LOWER(mc.name) = LOWER(?)", name). // Appetizer or appetizer by name params
		Where("products.is_available = ?", true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Category").
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productsRepo) IsProductExists(ctx context.Context, name string) (*entity.Product, error) {
	var prod entity.Product
	err := r.DB.WithContext(ctx).
		Where("LOWER(name) = LOWER(?) AND is_available = true", name).
		First(&prod).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &prod, err
}

// product yang aktif
func (r *productsRepo) HasActiveProducts(ctx context.Context, categoryID int64) (bool, error) {
	var count int64

	err := r.DB.WithContext(ctx).
		Model(&entity.Product{}).
		Where("category_id = ? AND is_available = true", categoryID).
		Count(&count).Error

	return count > 0, err
}
