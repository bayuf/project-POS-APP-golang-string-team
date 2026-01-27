package data

import (
	"fmt"

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
		entity.SeedUsers(),
		seedCategories,
		seedProducts,
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
		if err := db.FirstOrCreate(&cat, "name = ?", cat.Name).Error; err != nil {
			logger.Error("failed to seed menu category",
				zap.String("name", cat.Name),
				zap.Error(err),
			)
			return err
		}
		logger.Info("menu category ensured", zap.String("name", cat.Name))
	}

	return nil
}

func seedProducts(db *gorm.DB, logger *zap.Logger) error {
	// Get categories
	var categories []entity.MenuCategory
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	categoryMap := make(map[string]int64)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
func dataSeeds() []any {
	return []any{
		entity.SeedUsers(),
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
		var existing entity.Product
		err := db.
			Where("name = ? AND category_id = ?", prod.Name, prod.CategoryID).
			First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&prod).Error; err != nil {
				logger.Error("failed to seed product", zap.String("name", prod.Name), zap.Error(err))
				return err
			}
			logger.Info("product seeded", zap.String("name", prod.Name))
		} else if err != nil {
			return err
		} else {
			logger.Info("product already exists", zap.String("name", prod.Name))
		}
	}

	return nil
}
