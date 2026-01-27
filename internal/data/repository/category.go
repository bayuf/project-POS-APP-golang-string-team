package repository

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(page, limit int) ([]entity.MenuCategory, int64, error)
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
