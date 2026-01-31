package data

import (
	"fmt"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SeederFunc func(db *gorm.DB, logger *zap.Logger) error

func SeedAll(db *gorm.DB, logger *zap.Logger) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, seed := range dataSeeds() {
			if err := seed(tx, logger); err != nil {
				return fmt.Errorf("database seeding failed: %w", err)
			}
		}
		return nil
	})
}

func dataSeeds() []SeederFunc {
	return []SeederFunc{
		// entity.SeedUsers(),
		seedCategories,
		seedProducts,
		seedInventories,
		SeedTables,
		SeedReservations,
	}
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
		{CategoryID: categoryMap["Appetizer"], Name: "Garlic Bread", Price: 22000, IsAvailable: true},
		{CategoryID: categoryMap["Appetizer"], Name: "Chicken Spring Roll", Price: 28000, IsAvailable: true},

		{CategoryID: categoryMap["Main Course"], Name: "Nasi Goreng Spesial", Price: 45000, IsAvailable: true},
		{CategoryID: categoryMap["Main Course"], Name: "Mie Goreng Seafood", Price: 48000, IsAvailable: true},
		{CategoryID: categoryMap["Main Course"], Name: "Grilled Chicken Steak", Price: 65000, IsAvailable: true},

		{CategoryID: categoryMap["Dessert"], Name: "Chocolate Lava Cake", Price: 32000, IsAvailable: true},
		{CategoryID: categoryMap["Dessert"], Name: "Vanilla Ice Cream", Price: 25000, IsAvailable: true},

		{CategoryID: categoryMap["Beverage"], Name: "Hot Cappuccino", Price: 25000, IsAvailable: true},
		{CategoryID: categoryMap["Beverage"], Name: "Iced Lemon Tea", Price: 18000, IsAvailable: true},
		{CategoryID: categoryMap["Beverage"], Name: "Mineral Water", Price: 12000, IsAvailable: true},

		{CategoryID: categoryMap["Snack"], Name: "French Fries", Price: 20000, IsAvailable: true},
		{CategoryID: categoryMap["Snack"], Name: "Onion Rings", Price: 22000, IsAvailable: true},
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
