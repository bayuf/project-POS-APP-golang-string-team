package data

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
	// Example
	// &entity.User{},
	// &entity.Product{},
	)
}
