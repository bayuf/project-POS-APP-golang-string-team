package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderItem struct {
	ID        int64           `gorm:"primaryKey;autoIncrement"`
	OrderID   uuid.UUID       `gorm:"type:uuid;not null;index:idx_order_items_order_id"`
	ProductID int64           `gorm:"not null"`
	Quantity  int             `gorm:"not null"`
	Price     decimal.Decimal `gorm:"type:numeric(15,2);not null"`

	Order   Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID"`

	// composite unique
	// gorm tidak punya tag langsung, biasanya via migration manual
}
