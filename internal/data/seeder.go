package data

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		seeds := dataSeeds()
		for i := range seeds {
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(seeds[i]).Error; err != nil {
				name := reflect.TypeOf(seeds[i]).String()
				errMessage := err.Error()
				return fmt.Errorf("failed to seed %s: %s", name, errMessage)
			}
		}
		return nil
	})
}

func dataSeeds() []any {
	return []any{
		// Example
		// entity.SeedUsers(),
		// entity.SeedProducts(),
	}
}
