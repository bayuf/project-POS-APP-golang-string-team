package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepository UserRepositoryIface
	AuthRepository AuthRepositoryIface
}

func NewRepository(db *gorm.DB, logger *zap.Logger) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db, logger),
		AuthRepository: NewAuthRepository(db, logger),
	}
}
