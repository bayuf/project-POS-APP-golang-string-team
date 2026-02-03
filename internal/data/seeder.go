package data

import (
	"fmt"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SeederFunc func(db *gorm.DB, logger *zap.Logger) error

func SeedAll(db *gorm.DB, logger *zap.Logger) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, seed := range dataSeeds() {
			if err := seed(tx, logger); err != nil {
				logger.Error("seeding failed", zap.Error(err))
				return fmt.Errorf("database seeding failed: %w", err)
			}
		}
		return nil
	})
}

func dataSeeds() []SeederFunc {
	return []SeederFunc{
		seedUsers,
		seedCategories,
		seedProducts,
		seedInventories,
		seedPaymentMethods,
		SeedTables,
		SeedReservations,
		seedNotifications,
	}
}

func seedUsers(db *gorm.DB, logger *zap.Logger) error {
	now := time.Now()
	password, err := utils.HashString("admin12345")
	if err != nil {
		panic(err)
	}
	users := []entity.User{
		{
			ID:           uuid.New(),
			Name:         "Super Admin",
			Email:        "super@admin.com",
			PasswordHash: password,
			Phone:        "089111111111",
			Role:         "superadmin",
			Permissions: map[string]bool{
				entity.PermissionDashboard: true,
				entity.PermissionReports:   true,
				entity.PermissionInventory: true,
				entity.PermissionOrders:    true,
				entity.PermissionCustomers: true,
				entity.PermissionSettings:  true,
			},
			ShiftStart: "09:00:00",
			ShiftEnd:   "16:00:00",
			Address:    "Kab. Gresik",
			BirthDate:  time.Date(2000, 5, 19, 0, 0, 0, 0, time.UTC),
			Salary:     10000000,
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:           uuid.New(),
			Name:         "Bayu Firmansyah",
			Email:        "bayu19fr@gmail.com",
			PasswordHash: password,
			Phone:        "089111111111",
			Permissions: map[string]bool{
				entity.PermissionDashboard: true,
				entity.PermissionReports:   true,
				entity.PermissionInventory: true,
				entity.PermissionOrders:    true,
				entity.PermissionCustomers: true,
				entity.PermissionSettings:  false,
			},
			ShiftStart: "09:00:00",
			ShiftEnd:   "16:00:00",
			Role:       "admin",
			Address:    "Kab. Gresik",
			BirthDate:  time.Date(2000, 5, 19, 0, 0, 0, 0, time.UTC),
			Salary:     10000000,
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:           uuid.New(),
			Name:         "Alif Dwi Rahman",
			Email:        "alifdwirahman.alf@gmail.com",
			PasswordHash: password,
			Phone:        "089111111111",
			Role:         "admin",
			Permissions: map[string]bool{
				entity.PermissionDashboard: true,
				entity.PermissionReports:   true,
				entity.PermissionInventory: true,
				entity.PermissionOrders:    true,
				entity.PermissionCustomers: true,
				entity.PermissionSettings:  false,
			},
			ShiftStart: "09:00:00",
			ShiftEnd:   "16:00:00",
			Address:    "Kota Jakarta",
			BirthDate:  time.Date(2000, 5, 19, 0, 0, 0, 0, time.UTC),
			Salary:     10000000,
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&user, "name = ?", user.Name).FirstOrCreate(&user).Error; err != nil {
			logger.Error("failed to seed user",
				zap.String("name", user.Name),
				zap.Error(err),
			)
			return err
		}
		logger.Info("user ensured", zap.String("name", user.Name))
	}

	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count == 0 {
		return fmt.Errorf("category seeding failed, no data inserted")
	}

	return nil
}

func seedCategories(db *gorm.DB, logger *zap.Logger) error {
	categories := []entity.MenuCategory{
		{Name: "Appetizer"},
		{Name: "Main Course"},
		{Name: "Dessert"},
		{Name: "Beverage"},
		{Name: "Snack"},
	}

	for _, cat := range categories {
		if err := db.FirstOrCreate(&cat, "name = ?", cat.Name).FirstOrCreate(&cat).Error; err != nil {
			logger.Error("failed to seed menu category",
				zap.String("name", cat.Name),
				zap.Error(err),
			)
			return err
		}
		logger.Info("menu category ensured", zap.String("name", cat.Name))
	}

	var count int64
	db.Model(&entity.MenuCategory{}).Count(&count)
	if count == 0 {
		return fmt.Errorf("category seeding failed, no data inserted")
	}

	return nil
}

func seedProducts(db *gorm.DB, logger *zap.Logger) error {
	var categories []entity.MenuCategory
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	if len(categories) == 0 {
		return fmt.Errorf("no categories found, seedCategories must run first")
	}

	categoryMap := make(map[string]int64)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	products := []entity.Product{
		{CategoryID: categoryMap["Appetizer"], Name: "Garlic Bread", Price: decimal.New(22000, 0), IsAvailable: true},
		{CategoryID: categoryMap["Appetizer"], Name: "Chicken Spring Roll", Price: decimal.New(28000, 0), IsAvailable: true},

		{CategoryID: categoryMap["Main Course"], Name: "Nasi Goreng Spesial", Price: decimal.New(45000, 0), IsAvailable: true},
		{CategoryID: categoryMap["Main Course"], Name: "Mie Goreng Seafood", Price: decimal.New(48000, 0), IsAvailable: true},
		{CategoryID: categoryMap["Main Course"], Name: "Grilled Chicken Steak", Price: decimal.New(65000, 0), IsAvailable: true},

		{CategoryID: categoryMap["Dessert"], Name: "Chocolate Lava Cake", Price: decimal.New(32000, 0), IsAvailable: true},
		{CategoryID: categoryMap["Dessert"], Name: "Vanilla Ice Cream", Price: decimal.New(25000, 0), IsAvailable: true},

		{CategoryID: categoryMap["Beverage"], Name: "Hot Cappuccino", Price: decimal.New(25000, 0), IsAvailable: true},
		{CategoryID: categoryMap["Beverage"], Name: "Iced Lemon Tea", Price: decimal.New(18000, 0), IsAvailable: true},
		{CategoryID: categoryMap["Beverage"], Name: "Mineral Water", Price: decimal.New(12000, 0), IsAvailable: true},

		{CategoryID: categoryMap["Snack"], Name: "French Fries", Price: decimal.New(20000, 0), IsAvailable: true},
		{CategoryID: categoryMap["Snack"], Name: "Onion Rings", Price: decimal.New(22000, 0), IsAvailable: true},
	}

	for _, prod := range products {
		if prod.CategoryID == 0 {
			return fmt.Errorf("invalid category_id for product %s", prod.Name)
		}

		if err := db.
			Where("name = ? AND category_id = ?", prod.Name, prod.CategoryID).
			FirstOrCreate(&prod).Error; err != nil {

			logger.Error("failed to seed product",
				zap.String("name", prod.Name),
				zap.Error(err),
			)
			return err
		}

		logger.Info("product ensured",
			zap.Int64("id", prod.ID),
			zap.String("name", prod.Name),
		)
	}

	return nil
}

func seedInventories(db *gorm.DB, logger *zap.Logger) error {
	var products []entity.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	if len(products) == 0 {
		return fmt.Errorf("no products found, inventory seed aborted")
	}

	logger.Info("seeding inventories", zap.Int("products", len(products)))

	for _, p := range products {
		var existing entity.Inventory
		err := db.Unscoped().Where("product_id = ?", p.ID).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			inv := entity.Inventory{
				ProductID: p.ID,
				Stock:     15,
				Unit:      "pcs",
			}
			if err := db.Create(&inv).Error; err != nil {
				return err
			}
		} else if err == nil && existing.DeletedAt.Valid {
			db.Unscoped().Model(&existing).Update("deleted_at", nil)
			logger.Info("inventory restored", zap.Int64("product_id", p.ID))
		}
	}
	logger.Info("Inventories ensured")
	return nil
}

func SeedTables(db *gorm.DB, logger *zap.Logger) error {
	tables := []entity.RestaurantTable{
		{TableNumber: 1, Capacity: 2, IsActive: true},
		{TableNumber: 2, Capacity: 4, IsActive: true},
		{TableNumber: 3, Capacity: 4, IsActive: true},
		{TableNumber: 4, Capacity: 6, IsActive: true},
		{TableNumber: 5, Capacity: 8, IsActive: true},
	}

	for _, t := range tables {
		if err := db.Where("table_number = ?", t.TableNumber).FirstOrCreate(&t).Error; err != nil {
			return err
		}
	}
	logger.Info("restaurant tables ensured")
	return nil
}

func SeedReservations(db *gorm.DB, logger *zap.Logger) error {
	var table entity.RestaurantTable
	if err := db.Where("table_number = ?", 3).First(&table).Error; err != nil {
		logger.Warn("Table #3 not found, skipping reservation seed")
		return nil
	}

	layout := "2006-01-02 15:04"
	resTime1, _ := time.Parse(layout, "2026-02-10 19:00")
	resTime2, _ := time.Parse(layout, "2026-02-11 13:00")

	reservations := []entity.Reservation{
		{CustomerName: "John Doe", TableID: table.ID, ReservationTime: resTime1},
		{CustomerName: "Jane Smith", TableID: table.ID, ReservationTime: resTime2},
	}

	for _, r := range reservations {
		if err := db.Where("customer_name = ? AND reservation_time = ?", r.CustomerName, r.ReservationTime).
			FirstOrCreate(&r).Error; err != nil {
			return err
		}
	}
	logger.Info("reservations ensured")
	return nil
}

func seedNotifications(db *gorm.DB, logger *zap.Logger) error {
	var superAdmin entity.User
	// superadmin target notif
	if err := db.Where("role = ?", "superadmin").First(&superAdmin).Error; err != nil {
		logger.Warn("Superadmin not found, skipping notification seed")
		return nil
	}

	notifications := []entity.Notification{
		{
			ID:      uuid.New(),
			UserID:  superAdmin.ID,
			Title:   "Selamat Datang!",
			Message: "Sistem POS String Team sudah siap dipakai. Yuk, cek stok inventori kamu hari ini.",
			Status:  "new",
		},
		{
			ID:      uuid.New(),
			UserID:  superAdmin.ID,
			Title:   "Stok Hampir Habis",
			Message: "Produk 'Garlic Bread' tinggal kurang dari 5 pcs. Jangan lupa segera restock, ya.",
			Status:  "new",
		},
		{
			ID:      uuid.New(),
			UserID:  superAdmin.ID,
			Title:   "Reservasi Baru Masuk",
			Message: "Pelanggan atas nama John Doe telah memesan Meja #3 untuk pukul 19:00.",
			Status:  "read",
		},
	}

	for _, n := range notifications {
		if err := db.Where("title = ? AND message = ? AND user_id = ?", n.Title, n.Message, n.UserID).
			FirstOrCreate(&n).Error; err != nil {
			logger.Error("failed to seed notification", zap.Error(err))
			return err
		}
	}

	logger.Info("notifications ensured")
	return nil
}

func seedPaymentMethods(db *gorm.DB, logger *zap.Logger) error {
	paymentMethods := []entity.PaymentMethod{
		{Name: "Cash"},
		{Name: "Credit Card"},
		{Name: "Debit Card"},
	}

	if len(paymentMethods) == 0 {
		return fmt.Errorf("no payment methods found, payment method seed aborted")
	}

	logger.Info("seeding payment methods", zap.Int("payment methods", len(paymentMethods)))

	for _, p := range paymentMethods {
		pm := entity.PaymentMethod{
			Name: p.Name,
		}

		if err := db.
			Where("name = ?", p.Name).
			FirstOrCreate(&pm).Error; err != nil {
			return err
		}
	}

	return nil
}
