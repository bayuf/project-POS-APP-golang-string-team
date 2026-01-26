package repository

import (
	"context"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepositoryIface interface {
	CreateUser(ctx context.Context, newUser entity.User) error
	GetUserByID(ctx context.Context, ID uuid.UUID) (*entity.User, error)
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
	if err := r.db.WithContext(ctx).Create(&newUser).Error; err != nil {
		r.logger.Error("failed to create user", zap.Error(err))
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, ID).Error; err != nil {
		r.logger.Error("failed to get user by id", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserByID(ctx context.Context, ID uuid.UUID, updatedUser entity.User) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.User{}).Where("id = ?", ID).
		Updates(updatedUser).Error; err != nil {
		r.logger.Error("failed to update user by id", zap.Error(err))
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Delete(&entity.User{}, ID).Error; err != nil {
		r.logger.Error("failed to delete user by id", zap.Error(err))
		return err
	}
	return nil
}
