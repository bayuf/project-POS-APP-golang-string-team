package data

import (
	"strings"

	"gorm.io/gorm"
)

func runConstraints(db *gorm.DB) error {
	sqls := []string{
		// role
		`
		ALTER TABLE users
		ADD CONSTRAINT chk_users_role
		CHECK (role IN ('superadmin','admin','staff'));
		`,
		// price product
		`
		ALTER TABLE products
		ADD CONSTRAINT chk_products_price
		CHECK (price >= 0);
		`,
		// stock inventory
		`
		ALTER TABLE inventories
		ADD CONSTRAINT chk_inventories_stock
		CHECK (stock >= 0);
		`,
		// capacity table
		`
		ALTER TABLE restaurant_tables
		ADD CONSTRAINT chk_restaurant_tables_capacity
		CHECK (capacity > 0);
		`,
		// tax
		`
		ALTER TABLE orders
		ADD CONSTRAINT chk_orders_tax
		CHECK (tax >= 0);
		`,
		// total price in order
		`
		ALTER TABLE orders
		ADD CONSTRAINT chk_orders_total_price
		CHECK (total_price >= 0);
		`,
		// status in order
		`
		ALTER TABLE orders
		ADD CONSTRAINT chk_orders_status
		CHECK (status IN ('pending','paid','cancelled'));
		`,
		// quantity in order
		`
		ALTER TABLE order_items
		ADD CONSTRAINT chk_order_items_quantity
		CHECK (quantity > 0);
		`,
		// price in order item
		`
		ALTER TABLE order_items
		ADD CONSTRAINT chk_order_items_price
		CHECK (price >= 0);
		`,
		// orderid and productid must unique
		`
		ALTER TABLE order_items
		ADD CONSTRAINT uq_order_items_order_product
		UNIQUE (order_id, product_id);
		`,
		// status in notification
		`
		ALTER TABLE notifications
		ADD CONSTRAINT chk_notifications_status
		CHECK (status IN ('new','read'));
		`,
	}

	for _, q := range sqls {
		if err := db.Exec(q).Error; err != nil {
			// Abaikan jika constraint sudah ada
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
	}

	return nil
}
