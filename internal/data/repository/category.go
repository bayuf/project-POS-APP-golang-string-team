package repository

import (
	"context"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctg *entity.MenuCategory) error
	FindAll(page, limit int) ([]entity.MenuCategory, int64, error)
	FindById(ctx context.Context, ID int64) (*entity.MenuCategory, error)
	UpdateCategoryId(ctx context.Context, ID int64, ctg *entity.MenuCategory) (*entity.MenuCategory, error)

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

func (r *categoryRepo) Create(ctg *entity.MenuCategory) error {
	return r.DB.Create(ctg).Error
}

func (r *categoryRepo) FindAll(page, limit int) ([]entity.MenuCategory, int64, error) {
	var categories []entity.MenuCategory
	var total int64

	offset := (page - 1) * limit

	if err := r.DB.Model(&entity.MenuCategory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepo) FindById(ctx context.Context, id int64) (*entity.MenuCategory, error) {
	var ctg entity.MenuCategory
	if err := r.DB.WithContext(ctx).First(&ctg, id).Error; err != nil {
		r.Log.Error("failed to fetched category by id", zap.Error(err))
		return nil, err
	}

	return &ctg, nil
}

func (r *categoryRepo) UpdateCategoryId(ctx context.Context, id int64, ctg *entity.MenuCategory) (*entity.MenuCategory, error) {
	if err := r.DB.WithContext(ctx).Model(&entity.MenuCategory{}).Where("id = ?", id).Updates(ctg).Error; err != nil {

		r.Log.Error("failed to update category", zap.Error(err))
		return nil, err
	}

	return r.FindById(ctx, id)
}

func (r *categoryRepo) IsUniqueName(ctx context.Context, name string) (*entity.MenuCategory, error) {
	var category entity.MenuCategory

	err := r.DB.WithContext(ctx).
		Where("name = ?", name).
		First(&category).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}
