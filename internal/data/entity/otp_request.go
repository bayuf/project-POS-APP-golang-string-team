package entity

import (
	"time"

	"github.com/google/uuid"
)

type OTPRequest struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	OTPHash   string    `gorm:"type:text;not null"`
	ExpiredAt time.Time `gorm:"not null"`
	IsUsed    bool      `gorm:"default:false"`
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
