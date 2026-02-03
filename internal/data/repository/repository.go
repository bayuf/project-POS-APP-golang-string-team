package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	CategoryRepo    CategoryRepository
	ProductRepo     ProductsRepository
	UserRepository  UserRepositoryIface
	AuthRepository  AuthRepositoryIface
	InventoryRepo   InventoryRepository
	OrderRepository OrderRepositoryIface
	ReservationRepo ReservationRepository
	NotifRepo       NotificationRepository
}

func NewRepository(db *gorm.DB, logger *zap.Logger) *Repository {
	return &Repository{
		ProductRepo:     NewProductsRepository(db, logger),
		CategoryRepo:    NewCategoryRepository(db, logger),
		UserRepository:  NewUserRepository(db, logger),
		AuthRepository:  NewAuthRepository(db, logger),
		InventoryRepo:   NewInventoryRepository(db, logger),
		OrderRepository: NewOrderRepository(db, logger),
		ReservationRepo: NewReservationRepository(db, logger),
		NotifRepo:       NewNotificationRepository(db, logger),
	}
}
