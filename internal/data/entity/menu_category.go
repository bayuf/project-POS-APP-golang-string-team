package entity

import (
	"time"

	"gorm.io/gorm"
)

type MenuCategory struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:100;not null"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Products []Product `gorm:"foreignKey:CategoryID;references:ID"`
}
