package usecase

import (
	"context"

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
			ProgressStatus: "in proccess",
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
