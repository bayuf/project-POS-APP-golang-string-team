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

type InventoryService struct {
	repo   repository.InventoryRepository
	logger *zap.Logger
	tx     *gorm.DB
}

func NewInventoryService(repo repository.InventoryRepository, logg *zap.Logger, tx *gorm.DB) *InventoryService {
	return &InventoryService{
		repo:   repo,
		logger: logg,
		tx:     tx,
	}
}

func (u *InventoryService) CreateInventory(ctx context.Context, req dto.CreateInventory) (*entity.Inventory, error) {
	inventory := entity.Inventory{
		ProductID: req.ProductID,
		Stock:     req.Stock,
		Unit:      req.Unit,
	}

	return u.repo.CreateInventory(ctx, &inventory)
}

func (u *InventoryService) FindAllInventory(ctx context.Context, page, limit int) ([]entity.Inventory, dto.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	inventories, total, err := u.repo.FindAllInventory(ctx, page, limit)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   totalPage,
		TotalRecords: total,
	}

	return inventories, pagination, nil
}

func (u *InventoryService) GetDetailInventory(ctx context.Context, id int64) (*entity.Inventory, error) {
	inven, err := u.repo.FindByIdInventory(ctx, id)
	if err != nil {
		return nil, err
	}

	if inven == nil {
		return nil, errors.New("inventory not found")
	}

	return inven, nil
}

func (u *InventoryService) UpdateInventoryId(ctx context.Context, id int64, req dto.UpdateInventory) (*entity.Inventory, error) {
	checkExist, err := u.repo.FindByIdInventory(ctx, id)
	if err != nil {
		return nil, err
	}
	if checkExist == nil {
		return nil, errors.New("inventory not found")
	}

	updateInven := &entity.Inventory{
		Stock: req.Stock,
		Unit:  req.Unit,
	}

	return u.repo.UpdateInventoryId(ctx, id, updateInven)
}

func (u *InventoryService) DeleteInventoryId(ctx context.Context, id int64) error {
	exist, err := u.repo.FindByIdInventory(ctx, id)
	if err != nil {
		return err
	}

	if exist == nil {
		return errors.New("inventory not found")
	}

	return u.repo.DeletInventoryId(ctx, id)
}

func (u *InventoryService) SearchInventory(ctx context.Context, q dto.SearchInventoryQuery) ([]entity.Inventory, dto.Pagination, error) {

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}

	offset := (q.Page - 1) * q.Limit

	data, total, err := u.repo.SearchInventories(ctx, q, offset)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(q.Limit)))

	pagination := dto.Pagination{
		CurrentPage:  q.Page,
		Limit:        q.Limit,
		TotalPages:   totalPage,
		TotalRecords: total,
	}

	return data, pagination, nil
}
