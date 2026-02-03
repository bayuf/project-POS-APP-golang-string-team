package entity

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerName    string     `gorm:"column:customer_name;size:100"`
	TableID         int64      `gorm:"column:table_id;not null"`
	ReservationTime time.Time  `gorm:"column:reservation_time;not null;index:idx_reservations_time"`
	IsCancelled     bool       `gorm:"column:is_cancelled;default:false"`
	CreatedAt       time.Time  `gorm:"column:created_at;default:now()"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`

	Table RestaurantTable `gorm:"foreignKey:TableID"`
}
