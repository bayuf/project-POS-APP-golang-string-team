package usecase

import (
	"context"
	"errors"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderService struct {
	repo repository.OrderRepositoryIface
	log  *zap.Logger
	tx   *gorm.DB
}

func NewOrderService(repo repository.OrderRepositoryIface, log *zap.Logger, tx *gorm.DB) *OrderService {
	return &OrderService{
		repo: repo,
		log:  log,
		tx:   tx,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order dto.Order) (*dto.OrderResponse, error) {
	var IDs []int64
	for _, item := range order.Orders {
		IDs = append(IDs, item.ID)
	}

	// get item info
	items, err := s.repo.GetProductsInfoByID(ctx, IDs)
	if err != nil {
		return nil, err
	}

	orderID := uuid.New()
	subTotal := decimal.Zero
	if err := s.tx.Transaction(func(tx *gorm.DB) error {
		// generate order number
		orderNum := utils.GenerateOrderNum()

		// calculte price
		var totalPrice decimal.Decimal
		var totalPricePerItem []decimal.Decimal
		for i, item := range *items {
			// per item
			totalPricePerItem = append(totalPricePerItem, item.Price.Mul(decimal.NewFromInt(int64(order.Orders[i].Quantity))))

			// all item
			totalPrice = totalPrice.Add(item.Price.Mul(decimal.NewFromInt(int64(order.Orders[i].Quantity))))
		}

		// subtotal
		subTotal = totalPrice

		// count total price
		taxAmount := totalPrice.Mul(decimal.NewFromFloat(0.10))
		totalPrice = totalPrice.Add(taxAmount)

		// add order
		if err := s.repo.AddOrder(tx, ctx, &entity.Order{
			ID:             orderID,
			OrderNumber:    orderNum,
			CustomerName:   order.CustomerName,
			TableID:        &order.TableID,
			Tax:            taxAmount,
			TotalPrice:     totalPrice,
			ProgressStatus: "in the kitchen",
		}); err != nil {
			return err
		}

		// add item order
		itemsAdd := make([]entity.OrderItem, len(order.Orders))
		for i, item := range order.Orders {
			itemsAdd[i] = entity.OrderItem{
				OrderID:   orderID,
				ProductID: item.ID,
				Quantity:  item.Quantity,
				Price:     totalPricePerItem[i],
			}
		}
		if err := s.repo.AddOrderItems(tx, ctx, itemsAdd); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// get order detail
	orderDetail, err := s.repo.GetOrderDetailByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	orderDetailRes := dto.OrderResponse{
		OrderID:      orderDetail.ID,
		OrderNumber:  orderDetail.OrderNumber,
		CustomerName: orderDetail.CustomerName,
		TableID:      *orderDetail.TableID,
		Status:       orderDetail.OrderStatus,
		SubTotal:     subTotal,
		Tax:          orderDetail.Tax,
		TotalPrice:   orderDetail.TotalPrice,
	}

	return &orderDetailRes, nil
}

func (s *OrderService) PayOrder(ctx context.Context, paymentMethodID int64, orderID uuid.UUID) (*dto.OrderResponse, error) {
	// check paymentMethod
	paymentMethod, err := s.repo.GetPaymentMethodsById(ctx, paymentMethodID)
	if err != nil {
		return nil, err
	}

	if paymentMethod == nil {
		return nil, errors.New("payment method not found")
	}

	// update order
	if err := s.tx.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.EditOrder(tx, ctx, entity.Order{
			ID:              orderID,
			PaymentMethodID: &paymentMethodID,
			ProgressStatus:  "cooking now",
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// get order detail
	orderDetail, err := s.repo.GetOrderDetailByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// get sub total
	tax := orderDetail.TotalPrice.Mul(decimal.NewFromFloat(10.0)).Div(decimal.NewFromFloat(110.0))
	net := orderDetail.TotalPrice.Sub(tax)

	orderDetailRes := dto.OrderResponse{
		OrderID:      orderDetail.ID,
		OrderNumber:  orderDetail.OrderNumber,
		CustomerName: orderDetail.CustomerName,
		TableID:      *orderDetail.TableID,
		Status:       orderDetail.OrderStatus,
		SubTotal:     net,
		Tax:          orderDetail.Tax,
		TotalPrice:   orderDetail.TotalPrice,
	}

	return &orderDetailRes, nil
}

func (s *OrderService) EditOrder(ctx context.Context, newOrderData dto.Order, orderID uuid.UUID) error {

	var IDs []int64
	for _, item := range newOrderData.Orders {
		IDs = append(IDs, item.ID)
	}

	// get item info
	items, err := s.repo.GetProductsInfoByID(ctx, IDs)
	if err != nil {
		return err
	}

	// cek order
	orderData, err := s.repo.GetOrderDetailByID(ctx, orderID)
	if err != nil {
		return err
	}

	if orderData.ProgressStatus != "in the kitchen" {
		return errors.New("cant edit. order is already processed")
	}

	if err := s.tx.Transaction(func(tx *gorm.DB) error {
		// calculte price
		var totalPrice decimal.Decimal
		var totalPricePerItem []decimal.Decimal
		for i, item := range *items {
			// per item
			totalPricePerItem = append(totalPricePerItem, item.Price.Mul(decimal.NewFromInt(int64(newOrderData.Orders[i].Quantity))))

			// all item
			totalPrice = totalPrice.Add(item.Price.Mul(decimal.NewFromInt(int64(newOrderData.Orders[i].Quantity))))
		}

		// count total price
		taxAmount := totalPrice.Mul(decimal.NewFromFloat(0.10))
		totalPrice = totalPrice.Add(taxAmount)

		// update order
		if err := s.repo.EditOrder(tx, ctx, entity.Order{
			ID:           orderData.ID,
			CustomerName: newOrderData.CustomerName,
			Tax:          taxAmount,
			TotalPrice:   totalPrice,
		}); err != nil {
			return err
		}

		// update item order
		itemsAdd := make([]entity.OrderItem, len(newOrderData.Orders))
		for i, item := range newOrderData.Orders {
			itemsAdd[i] = entity.OrderItem{
				OrderID:   orderID,
				ProductID: item.ID,
				Quantity:  item.Quantity,
				Price:     totalPricePerItem[i],
			}
		}
		if err := s.repo.UpdateOrderItems(tx, ctx, itemsAdd); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	if err := s.repo.DeleteOrder(ctx, orderID); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) ProcessOrder(ctx context.Context, orderID uuid.UUID) error {
	// get order detail
	orderItem, err := s.repo.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	products := make([]entity.OrderItem, len(*orderItem))
	for _, item := range *orderItem {
		products = append(products, item)
	}

	if err := s.tx.Transaction(func(tx *gorm.DB) error {

		if err := s.repo.UpdateStockInventories(tx, ctx, products); err != nil {
			return err
		}

		if err := s.repo.ProcessOrder(tx, ctx, orderID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) CompleteOrder(ctx context.Context, orderID uuid.UUID) error {
	if err := s.repo.CompleteOrder(ctx, orderID); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) GetPaymentMethods(ctx context.Context) ([]dto.PaymentMethodResponse, error) {
	method, err := s.repo.GetPaymentMethods(ctx)
	if err != nil {
		return nil, err
	}

	var paymentMethods []dto.PaymentMethodResponse
	for _, m := range method {
		paymentMethods = append(paymentMethods, dto.PaymentMethodResponse{
			ID:   int64(m.ID),
			Name: m.Name,
		})
	}

	return paymentMethods, nil
}

func (s *OrderService) GetTables(ctx context.Context) ([]entity.RestaurantTable, error) {
	tables, err := s.repo.GetTables(ctx)
	if err != nil {
		return nil, err
	}

	return tables, nil
}
