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
		inv := entity.Inventory{
			ProductID: p.ID,
			Stock:     15,
			Unit:      "pcs",
		}

		if err := db.
			Where("product_id = ?", p.ID).
			FirstOrCreate(&inv).Error; err != nil {
			return err
		}
	}

	return nil
}
