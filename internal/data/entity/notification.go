package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_notifications_user_id"`
	Title     string    `gorm:"size:150;not null"`
	Message   string    `gorm:"type:text;not null"`
	Status    string    `gorm:"size:10;default:new"` // new, read
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
