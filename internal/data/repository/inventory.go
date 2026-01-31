package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	FindAllInventory(ctx context.Context, page, limit int) ([]entity.Inventory, int64, error)
	CreateInventory(ctx context.Context, inven *entity.Inventory) (*entity.Inventory, error)
	FindByIdInventory(ctx context.Context, id int64) (*entity.Inventory, error)
	UpdateInventoryId(ctx context.Context, id int64, inven *entity.Inventory) (*entity.Inventory, error)
	DeletInventoryId(ctx context.Context, id int64) error

	SearchInventories(ctx context.Context, q dto.SearchInventoryQuery, offset int) ([]entity.Inventory, int64, error)
}

type InventoryRepo struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewInventoryRepository(db *gorm.DB, logg *zap.Logger) *InventoryRepo {
	return &InventoryRepo{
		DB:     db,
		Logger: logg,
	}
}

func (r *InventoryRepo) CreateInventory(ctx context.Context, inven *entity.Inventory) (*entity.Inventory, error) {

	if err := r.DB.WithContext(ctx).Create(inven).Error; err != nil {
		return nil, err
	}

	// with prod data
	if err := r.DB.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		First(inven, inven.ID).Error; err != nil {
		return nil, err
	}

	return inven, nil
}

func (r *InventoryRepo) FindAllInventory(ctx context.Context, page, limit int) ([]entity.Inventory, int64, error) {

	var inventories []entity.Inventory
	var total int64

	offset := (page - 1) * limit

	if err := r.DB.WithContext(ctx).
		Model(&entity.Inventory{}).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// retrieve + join
	if err := r.DB.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		Limit(limit).
		Offset(offset).
		Find(&inventories).Error; err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}

func (r *InventoryRepo) FindByIdInventory(ctx context.Context, id int64) (*entity.Inventory, error) {
	var inven entity.Inventory
	err := r.DB.WithContext(ctx).Preload("Product").First(&inven, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &inven, nil
}

func (r *InventoryRepo) UpdateInventoryId(ctx context.Context, id int64, inven *entity.Inventory) (*entity.Inventory, error) {
	if err := r.DB.WithContext(ctx).Model(&entity.Inventory{}).Where("id = ?", id).Updates(inven).Error; err != nil {

		r.Logger.Error("failed to update inventory", zap.Error(err))
		return nil, err
	}

	return r.FindByIdInventory(ctx, id)
}

func (r *InventoryRepo) DeletInventoryId(ctx context.Context, id int64) error {
	result := r.DB.WithContext(ctx).
		Delete(&entity.Inventory{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *InventoryRepo) SearchInventories(
	ctx context.Context,
	q dto.SearchInventoryQuery,
	offset int,
) ([]entity.Inventory, int64, error) {

	var inventories []entity.Inventory
	var total int64

	db := r.DB.WithContext(ctx).
		Model(&entity.Inventory{}).
		Joins("JOIN products ON products.id = inventories.product_id").
		Joins("JOIN menu_categories ON menu_categories.id = products.category_id")

	if q.ProductName != "" {
		db = db.Where(
			"LOWER(products.name) LIKE ?",
			"%"+strings.ToLower(q.ProductName)+"%",
		)
	}

	if q.CategoryID != nil {
		db = db.Where("products.category_id = ?", *q.CategoryID)
	}

	if q.MinStock != nil {
		db = db.Where("inventories.stock >= ?", *q.MinStock)
	}

	if q.MaxStock != nil {
		db = db.Where("inventories.stock <= ?", *q.MaxStock)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Preload("Product").
		Preload("Product.Category").
		Limit(q.Limit).
		Offset(offset).
		Find(&inventories).Error; err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}
