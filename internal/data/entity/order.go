package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderNumber     string    `gorm:"size:50;not null;unique"`
	CustomerName    string    `gorm:"size:100"`
	TableID         *int64
	PaymentMethodID *int64
	Tax             float64   `gorm:"type:numeric(5,2);not null"`
	TotalPrice      float64   `gorm:"type:numeric(14,2);not null"`
	Status          string    `gorm:"size:20;default:pending"` // pending, paid, cancelled
	CreatedAt       time.Time `gorm:"index:idx_orders_created_at"`
	UpdatedAt       *time.Time

	Table         *RestaurantTable
	PaymentMethod *PaymentMethod
	Items         []OrderItem
}
