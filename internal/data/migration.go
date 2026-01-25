package data

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.MenuCategory{},
		&entity.RestaurantTable{},
		&entity.PaymentMethod{},

		&entity.Product{},
		&entity.Inventory{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.Reservation{},
		&entity.OTPRequest{},
		&entity.Notification{},
	)
}
