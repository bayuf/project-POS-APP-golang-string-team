package repository

import (
	"context"
	"errors"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, ctg *entity.MenuCategory) error
	FindAll(ctx context.Context, page, limit int) ([]entity.MenuCategory, int64, error)
	FindById(ctx context.Context, ID int64) (*entity.MenuCategory, error)
	UpdateCategoryId(ctx context.Context, ID int64, ctg *entity.MenuCategory) (*entity.MenuCategory, error)
	DeleteCategoryId(ctx context.Context, id int64) error

	IsUniqueName(ctx context.Context, name string) (*entity.MenuCategory, error)
}

type categoryRepo struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewCategoryRepository(db *gorm.DB, log *zap.Logger) CategoryRepository {
	return &categoryRepo{
		DB:  db,
		Log: log,
	}
}

func (r *categoryRepo) Create(ctx context.Context, ctg *entity.MenuCategory) error {
	return r.DB.WithContext(ctx).Create(ctg).Error
}

func (r *categoryRepo) FindAll(ctx context.Context, page, limit int) ([]entity.MenuCategory, int64, error) {
	var categories []entity.MenuCategory
	var total int64

	offset := (page - 1) * limit

	// total
	if err := r.DB.WithContext(ctx).
		Model(&entity.MenuCategory{}).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// data + preload products
	if err := r.DB.WithContext(ctx).
		Preload("Products").
		Limit(limit).
		Offset(offset).
		Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepo) FindById(ctx context.Context, id int64) (*entity.MenuCategory, error) {
	var category entity.MenuCategory

	err := r.DB.Debug().
		Model(&entity.MenuCategory{}).
		Preload("Products").
		Where("id = ?", id).
		First(&category).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &category, nil
}

func (r *categoryRepo) UpdateCategoryId(ctx context.Context, id int64, ctg *entity.MenuCategory) (*entity.MenuCategory, error) {
	if err := r.DB.WithContext(ctx).Model(&entity.MenuCategory{}).Where("id = ?", id).Updates(ctg).Error; err != nil {

		r.Log.Error("failed to update category", zap.Error(err))
		return nil, err
	}

	return r.FindById(ctx, id)
}

func (r *categoryRepo) DeleteCategoryId(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		catRepo := NewCategoryRepository(tx, r.Log)
		prodRepo := NewProductsRepository(tx, r.Log)

		ctg, err := catRepo.FindById(ctx, id)
		if err != nil {
			return err
		}
		if ctg == nil {
			return errors.New("category not found")
		}

		// produk yang active
		hasProd, err := prodRepo.HasActiveProducts(ctx, id)
		if err != nil {
			return err
		}
		if hasProd {
			return errors.New("category still has active products")
		}

		now := time.Now()

		return tx.WithContext(ctx).
			Model(&entity.MenuCategory{}).
			Where("id = ?", id).
			Update("deleted_at", &now).
			Error
	})
}

func (r *categoryRepo) IsUniqueName(ctx context.Context, name string) (*entity.MenuCategory, error) {
	var category entity.MenuCategory

	err := r.DB.WithContext(ctx).
		Where("LOWER(name) = LOWER(?)", name).
		First(&category).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}
