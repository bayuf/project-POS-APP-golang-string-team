# Project POS APP Golang

## Database Migrate
### After auto migrate database add constraint (Opsional).

table users
```
ALTER TABLE users
ADD CONSTRAINT chk_users_role
CHECK (role IN ('superadmin','admin','staff'));
```

table products
```
ALTER TABLE products
ADD CONSTRAINT chk_products_price
CHECK (price >= 0);
```

table inventories
```
ALTER TABLE inventories
ADD CONSTRAINT chk_inventories_stock
CHECK (stock >= 0);
```

table Restaurant Tables
```
ALTER TABLE restaurant_tables
ADD CONSTRAINT chk_restaurant_tables_capacity
CHECK (capacity > 0);
```

table orders
```
ALTER TABLE orders
ADD CONSTRAINT chk_orders_tax
CHECK (tax >= 0);

ALTER TABLE orders
ADD CONSTRAINT chk_orders_total_price
CHECK (total_price >= 0);

ALTER TABLE orders
ADD CONSTRAINT chk_orders_status
CHECK (status IN ('pending','paid','cancelled'));
```

table order_items
```
ALTER TABLE order_items
ADD CONSTRAINT chk_order_items_quantity
CHECK (quantity > 0);

ALTER TABLE order_items
ADD CONSTRAINT chk_order_items_price
CHECK (price >= 0);

ALTER TABLE order_items
ADD CONSTRAINT uq_order_items_order_product
UNIQUE (order_id, product_id);
```

table notifications
```
ALTER TABLE notifications
ADD CONSTRAINT chk_notifications_status
CHECK (status IN ('new','read'));
```
