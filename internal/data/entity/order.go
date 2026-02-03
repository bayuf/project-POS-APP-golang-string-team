package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderNumber     string    `gorm:"size:50;not null;unique"`
	CustomerName    string    `gorm:"size:100"`
	TableID         *int64
	PaymentMethodID *int64
	Tax             decimal.Decimal `gorm:"type:numeric(15,2);not null"`
	TotalPrice      decimal.Decimal `gorm:"type:numeric(15,2);not null"`
	ProgressStatus  string          `gorm:"size:20;default:in proccess"`
	OrderStatus     string          `gorm:"size:20;default:in proccess"`
	CreatedAt       time.Time       `gorm:"index:idx_orders_created_at"`
	UpdatedAt       *time.Time

	Table         *RestaurantTable
	PaymentMethod *PaymentMethod
	Items         []OrderItem
}
