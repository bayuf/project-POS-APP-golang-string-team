package entity

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerName    string    `gorm:"size:100"`
	TableID         int64     `gorm:"not null"`
	ReservationTime time.Time `gorm:"not null;index:idx_reservations_time"`
	IsCancelled     bool      `gorm:"default:false"`
	CreatedAt       time.Time
	UpdatedAt       *time.Time

	Table RestaurantTable `gorm:"foreignKey:TableID"`
}
