package data

import (
	"gorm.io/gorm"
)

func addConstraintIfNotExists(db *gorm.DB, name string, sql string) error {
	var exists bool

	checkSQL := `
		SELECT EXISTS (
			SELECT 1
			FROM pg_constraint
			WHERE conname = ?
		);
	`

	if err := db.Raw(checkSQL, name).Scan(&exists).Error; err != nil {
		return err
	}

	if exists {
		return nil
	}

	return db.Exec(sql).Error
}

func runConstraints(db *gorm.DB) error {
	constraints := []struct {
		name string
		sql  string
	}{
		{
			"chk_users_role",
			`
			ALTER TABLE users
			ADD CONSTRAINT chk_users_role
			CHECK (role IN ('superadmin','admin','staff'));
			`,
		},
		{
			"chk_products_price",
			`
			ALTER TABLE products
			ADD CONSTRAINT chk_products_price
			CHECK (price >= 0);
			`,
		},
		{
			"chk_inventories_stock",
			`
			ALTER TABLE inventories
			ADD CONSTRAINT chk_inventories_stock
			CHECK (stock >= 0);
			`,
		},
		{
			"chk_restaurant_tables_capacity",
			`
			ALTER TABLE restaurant_tables
			ADD CONSTRAINT chk_restaurant_tables_capacity
			CHECK (capacity > 0);
			`,
		},
		{
			"chk_orders_tax",
			`
			ALTER TABLE orders
			ADD CONSTRAINT chk_orders_tax
			CHECK (tax >= 0);
			`,
		},
		{
			"chk_orders_total_price",
			`
			ALTER TABLE orders
			ADD CONSTRAINT chk_orders_total_price
			CHECK (total_price >= 0);
			`,
		},
		{
			"chk_progress_orders_status",
			`
			ALTER TABLE orders
			ADD CONSTRAINT chk_progress_orders_status
			CHECK (progress_status IN ('cancelled','in the kitchen','cooking now','ready to serve'));
			`,
		},
		{
			"chk_orders_status",
			`
			ALTER TABLE orders
			ADD CONSTRAINT chk_orders_status
			CHECK (order_status IN ('in proccess','ready','completed', 'cancelled'));
			`,
		},
		{
			"chk_order_items_quantity",
			`
			ALTER TABLE order_items
			ADD CONSTRAINT chk_order_items_quantity
			CHECK (quantity > 0);
			`,
		},
		{
			"chk_order_items_price",
			`
			ALTER TABLE order_items
			ADD CONSTRAINT chk_order_items_price
			CHECK (price >= 0);
			`,
		},
		{
			"uq_order_items_order_product",
			`
			ALTER TABLE order_items
			ADD CONSTRAINT uq_order_items_order_product
			UNIQUE (order_id, product_id);
			`,
		},
		{
			"chk_notifications_status",
			`
			ALTER TABLE notifications
			ADD CONSTRAINT chk_notifications_status
			CHECK (status IN ('new','read'));
			`,
		},
	}

	for _, c := range constraints {
		if err := addConstraintIfNotExists(db, c.name, c.sql); err != nil {
			return err
		}
	}

	return nil
}
