package data

import (
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
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
	); err != nil {
		return err
	}

	if err := manual(db); err != nil {
		return err
	}

	if err := runConstraints(db); err != nil {
		return err
	}

	return nil
}

func manual(db *gorm.DB) error {
	type columnInfo struct {
		ColumnName string
		DataType   string
	}

	var cols []columnInfo

	// cek tipe kolom di postgres
	err := db.Raw(`
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_name = 'users'
		  AND column_name IN ('shift_start', 'shift_end');
	`).Scan(&cols).Error
	if err != nil {
		return err
	}

	for _, col := range cols {
		if col.DataType != "time without time zone" {
			// paksa ubah ke TIME
			if err := db.Exec(
				`ALTER TABLE users
				 ALTER COLUMN ` + col.ColumnName + `
				 TYPE TIME
				 USING ` + col.ColumnName + `::TIME;`,
			).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
