package entity

type PaymentMethod struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"size:50;not null;unique"`
}
