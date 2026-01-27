package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID    uuid.UUID `gorm:"not null"`
	ExpiredAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	UsedAt    *time.Time

	CreatedAt time.Time `gorm:"not null"`

	User User `gorm:"foreignKey:UserID"`
}
