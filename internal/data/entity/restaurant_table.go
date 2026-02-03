package entity

type RestaurantTable struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	TableNumber int   `gorm:"not null;unique"`
	Capacity    int   `gorm:"not null"`
	IsActive    bool  `gorm:"default:true"`
}
