package entity

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	ProductID int64  `gorm:"not null;unique"`
	Stock     int    `gorm:"not null"`
	Unit      string `gorm:"size:20"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}
