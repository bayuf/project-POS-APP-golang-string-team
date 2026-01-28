package usecase

import (
	"context"
	"errors"
	"math"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"gorm.io/gorm"
)

func (u UseCase) CreateCategory(ctx context.Context, req dto.CreateCategory) (*entity.MenuCategory, error) {
	//  unique name
	existing, err := u.repo.CategoryRepo.IsUniqueName(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("category name already exists")
	}

	category := &entity.MenuCategory{
		Name: req.Name,
	}

	if err := u.repo.CategoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (u UseCase) FindAllCategories(page, limit int) ([]entity.MenuCategory, dto.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	categories, total, err := u.repo.CategoryRepo.FindAll(context.Background(), page, limit)
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

	return categories, pagination, nil
}

func (u UseCase) GetCategoryByID(ctx context.Context, id int64) (*entity.MenuCategory, error) {

	category, err := u.repo.CategoryRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("category not found")
	}

	return category, nil
}

func (u UseCase) UpdateCategoryId(ctx context.Context, id int64, req dto.CreateCategory) (*entity.MenuCategory, error) {

	existing, err := u.repo.CategoryRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("category not found")
	}

	// unique name check (exclude current id)
	duplicate, err := u.repo.CategoryRepo.IsUniqueName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if duplicate != nil && duplicate.ID != id {
		return nil, errors.New("category name already exists")
	}

	update := &entity.MenuCategory{
		Name: req.Name,
	}

	return u.repo.CategoryRepo.UpdateCategoryId(ctx, id, update)
}

func (u *UseCase) DeleteCategory(ctx context.Context, id int64) error {

	return u.tx.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		catRepo := repository.NewCategoryRepository(tx, u.log)
		prodRepo := repository.NewProductsRepository(tx, u.log)

		ctg, err := catRepo.FindById(ctx, id)
		if err != nil {
			return err
		}
		if ctg == nil {
			return errors.New("category not found")
		}

		hasProd, err := prodRepo.HasActiveProducts(ctx, id)
		if err != nil {
			return err
		}
		if hasProd {
			return errors.New("category still has active products")
		}

		return catRepo.DeleteCategoryId(ctx, id)
	})
}
