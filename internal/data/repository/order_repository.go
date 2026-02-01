package repository

import (
	"context"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepositoryIface interface {
	GetOrderDetailByID(ctx context.Context, ID uuid.UUID) (*entity.Order, error)
	GetProductsInfoByID(ctx context.Context, items []int64) (*[]entity.Product, error)
	AddOrder(tx *gorm.DB, ctx context.Context, order *entity.Order) error
	AddOrderItems(tx *gorm.DB, ctx context.Context, itemOrders []entity.OrderItem) error
	GetPaymentMethodsById(ctx context.Context, id int64) (*entity.PaymentMethod, error)
	UpdateOrder(tx *gorm.DB, ctx context.Context, order entity.Order) error
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

func (r *OrderRepository) UpdateOrder(tx *gorm.DB, ctx context.Context, order entity.Order) error {
	if err := tx.WithContext(ctx).
		Where("id = ?", order.ID).
		Updates(&order).
		Error; err != nil {
		r.logger.Error("failed to update order", zap.Error(err))
		return err
	}

	return nil
}
