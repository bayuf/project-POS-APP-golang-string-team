package repository

import (
	"context"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepositoryIface interface {
	GetTables(ctx context.Context) ([]entity.RestaurantTable, error)
	GetPaymentMethods(ctx context.Context) ([]entity.PaymentMethod, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID uuid.UUID) (*[]entity.OrderItem, error)
	GetOrderDetailByID(ctx context.Context, ID uuid.UUID) (*entity.Order, error)
	GetProductsInfoByID(ctx context.Context, items []int64) (*[]entity.Product, error)
	AddOrder(tx *gorm.DB, ctx context.Context, order *entity.Order) error
	AddOrderItems(tx *gorm.DB, ctx context.Context, itemOrders []entity.OrderItem) error
	GetPaymentMethodsById(ctx context.Context, id int64) (*entity.PaymentMethod, error)
	EditOrder(tx *gorm.DB, ctx context.Context, order entity.Order) error
	UpdateOrderItems(tx *gorm.DB, ctx context.Context, itemOrders []entity.OrderItem) error
	DeleteOrder(ctx context.Context, orderID uuid.UUID) error
	ProcessOrder(tx *gorm.DB, ctx context.Context, orderID uuid.UUID) error
	CompleteOrder(ctx context.Context, orderID uuid.UUID) error
	UpdateStockInventories(tx *gorm.DB, ctx context.Context, items []entity.OrderItem) error
}

type OrderRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewOrderRepository(db *gorm.DB, log *zap.Logger) *OrderRepository {
	return &OrderRepository{
		db:     db,
		logger: log,
	}
}

func (r *OrderRepository) GetTables(ctx context.Context) ([]entity.RestaurantTable, error) {
	var tables []entity.RestaurantTable
	now := time.Now()

	err := r.db.WithContext(ctx).
		Table("restaurant_tables").
		Joins(`
				LEFT JOIN reservations
				ON reservations.table_id = restaurant_tables.id
				AND reservations.is_cancelled = false
				AND reservations.reservation_time >= ?
			`, now).
		Where("restaurant_tables.is_active = ?", true).
		Where("reservations.id IS NULL").
		Find(&tables).
		Error

	if err != nil {
		r.logger.Error("failed to get ready tables", zap.Error(err))
		return nil, err
	}

	return tables, nil
}

func (r *OrderRepository) GetPaymentMethods(ctx context.Context) ([]entity.PaymentMethod, error) {
	var paymentMethods []entity.PaymentMethod

	if err := r.db.WithContext(ctx).
		Find(&paymentMethods).Error; err != nil {
		r.logger.Error("failed to get payment methods", zap.Error(err))
		return nil, err
	}

	return paymentMethods, nil
}

func (r *OrderRepository) GetPaymentMethodsById(ctx context.Context, id int64) (*entity.PaymentMethod, error) {
	var paymentMethods entity.PaymentMethod

	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Find(&paymentMethods).Error; err != nil {
		r.logger.Error("failed to get payment methods", zap.Error(err))
		return nil, err
	}

	return &paymentMethods, nil
}

func (r *OrderRepository) GetOrderDetailByID(ctx context.Context, ID uuid.UUID) (*entity.Order, error) {
	var order entity.Order

	if err := r.db.WithContext(ctx).First(&order, ID).
		Error; err != nil {
		r.logger.Error("failed to get order detail", zap.Error(err))
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetProductsInfoByID(ctx context.Context, items []int64) (*[]entity.Product, error) {
	var products []entity.Product

	if err := r.db.WithContext(ctx).
		Preload("Inventory").
		Joins("JOIN inventories ON products.id = inventories.product_id").
		Where("products.id IN ?", items).
		Where("products.is_available = ?", true).
		Where("inventories.stock > 0").
		Find(&products).
		Error; err != nil {
		r.logger.Error("failed to get products info by id", zap.Error(err))
		return nil, err
	}

	return &products, nil
}

func (r *OrderRepository) GetOrderItemsByOrderID(ctx context.Context, orderID uuid.UUID) (*[]entity.OrderItem, error) {
	var items []entity.OrderItem

	if err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Find(&items).
		Error; err != nil {
		r.logger.Error("failed to get order items by id", zap.Error(err))
		return nil, err
	}

	return &items, nil
}

func (r *OrderRepository) AddOrder(tx *gorm.DB, ctx context.Context, order *entity.Order) error {
	if err := r.db.WithContext(ctx).Create(&order).Error; err != nil {
		r.logger.Error("failed to add order", zap.Error(err))
		return err
	}

	return nil
}

func (r *OrderRepository) AddOrderItems(tx *gorm.DB, ctx context.Context, itemOrders []entity.OrderItem) error {
	if err := r.db.WithContext(ctx).Create(&itemOrders).Error; err != nil {
		r.logger.Error("failed to add order item", zap.Error(err))
		return err
	}

	return nil
}

func (r *OrderRepository) EditOrder(tx *gorm.DB, ctx context.Context, order entity.Order) error {
	err := tx.WithContext(ctx).
		Model(&entity.Order{}).
		Where("id = ?", order.ID).
		Where("progress_status = ?", "in the kitchen").
		Updates(&order).Error

	if err != nil {
		r.logger.Error("failed to update order", zap.Error(err))
		return err
	}

	return nil
}

func (r *OrderRepository) UpdateOrderItems(tx *gorm.DB, ctx context.Context, itemOrders []entity.OrderItem) error {
	for _, item := range itemOrders {
		if err := tx.WithContext(ctx).Model(&entity.OrderItem{}).
			Where("order_id = ?", item.OrderID).
			Where("product_id = ?", item.ProductID).
			Updates(item).Error; err != nil {
			r.logger.Error("failed to update order item", zap.Error(err))
			return err
		}
	}

	return nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.Order{}).
		Where("id = ?", orderID).
		Where("order_status = ?", "in proccess").
		Update("progress_status", "cancelled").
		Update("order_status", "cancelled").
		Error; err != nil {
		r.logger.Error("failed to delete order", zap.Error(err))
		return err
	}

	return nil
}

func (r *OrderRepository) ProcessOrder(tx *gorm.DB, ctx context.Context, orderID uuid.UUID) error {
	if err := tx.WithContext(ctx).
		Model(&entity.Order{}).
		Where("id = ?", orderID).
		Where("progress_status = ?", "cooking now").
		Updates(map[string]any{
			"progress_status": "ready to serve",
			"order_status":    "ready",
		}).
		Error; err != nil {
		r.logger.Error("failed to process order", zap.Error(err))
		return err
	}

	return nil
}

func (r *OrderRepository) CompleteOrder(ctx context.Context, orderID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.Order{}).
		Where("id = ?", orderID).
		Where("progress_status = ?", "ready to serve").
		Where("order_status = ?", "ready").
		Update("order_status", "completed").
		Error; err != nil {
		r.logger.Error("failed to complete order", zap.Error(err))
		return err
	}

	return nil
}

func (r *OrderRepository) UpdateStockInventories(tx *gorm.DB, ctx context.Context, items []entity.OrderItem) error {
	for _, item := range items {
		if err := tx.WithContext(ctx).
			Model(&entity.Inventory{}).
			Where("product_id = ?", item.ProductID).
			Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			r.logger.Error("failed to update stock inventory", zap.Error(err))
			return err
		}
	}

	return nil
}
