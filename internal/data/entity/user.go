package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:100;not null;unique"`
	PasswordHash string    `gorm:"type:text;not null;default:''"`
	Role         string    `gorm:"size:20;not null"` // superadmin, admin, staff
	IsActive     bool      `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	OTPRequests   []OTPRequest
	Notifications []Notification
}
