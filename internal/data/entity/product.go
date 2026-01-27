package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          int64   `gorm:"primaryKey;autoIncrement"`
	CategoryID  int64   `gorm:"not null;index:idx_products_category_id"`
	Name        string  `gorm:"size:150;not null"`
	Price       float64 `gorm:"type:numeric(12,2);not null"`
	IsAvailable bool    `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Category  *MenuCategory `gorm:"foreignKey:CategoryID;references:ID"`
	Inventory *Inventory
}
