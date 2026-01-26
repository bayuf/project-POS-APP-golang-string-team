package repository

import (
	"context"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepositoryIface interface {
	CreateUser(ctx context.Context, newUser entity.User) error
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
