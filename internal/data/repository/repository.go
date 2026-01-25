package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository(db *gorm.DB, logger *zap.Logger) *Repository {
	return &Repository{}
}
