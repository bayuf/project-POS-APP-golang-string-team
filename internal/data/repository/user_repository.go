package repository

import (
	"context"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepositoryIface interface {
	ListUser(ctx context.Context, f dto.UserFilterRequest) (*[]entity.User, int64, error)
	CreateUser(ctx context.Context, newUser entity.User) error
	GetUserByID(ctx context.Context, ID uuid.UUID) (*entity.User, error)
	UpdateUserByID(ctx context.Context, ID uuid.UUID, updatedUser entity.User) error
	DeleteUserByID(ctx context.Context, ID uuid.UUID) error
}

type UserRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, newUser entity.User) error {
	if err := r.db.WithContext(ctx).
		Create(&newUser).
		Error; err != nil {
		r.logger.Error("failed to create user", zap.Error(err))
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).
		First(&user, ID).
		Error; err != nil {
		r.logger.Error("failed to get user by id", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserByID(ctx context.Context, ID uuid.UUID, updatedUser entity.User) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", ID).
		Updates(updatedUser).
		Error; err != nil {
		r.logger.Error("failed to update user by id", zap.Error(err))
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ? AND role != ?", ID, "superadmin").
		Update("is_active", false).
		Error; err != nil {
		r.logger.Error("failed to deactivate user by id", zap.Error(err))
		return err
	}

	if err := r.db.WithContext(ctx).
		Delete(&entity.User{}, ID).
		Where("role != ?", "superadmin").
		Error; err != nil {
		r.logger.Error("failed to delete user by id", zap.Error(err))
		return err
	}

	return nil
}

func (r *UserRepository) ListUser(ctx context.Context, f dto.UserFilterRequest) (*[]entity.User, int64, error) {
	var users []entity.User
	var totalItems int64

	query := r.db.Model(&entity.User{})

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	switch f.SortBy {
	case "email":
		query = query.Order("email asc")
	case "name":
		query = query.Order("name asc")
	default:
		query = query.Order("created_at desc")
	}

	offset := (f.Page - 1) * f.Limit

	err := query.Limit(f.Limit).Offset(offset).Find(&users).Error

	return &users, totalItems, err
}
