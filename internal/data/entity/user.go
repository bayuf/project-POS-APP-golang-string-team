package entity

import (
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:100;not null;unique"`
	PasswordHash string    `gorm:"type:text;not null;default:''"`
	Role         string    `gorm:"size:20;not null"` // superadmin, admin, staff
	Address      string    `gorm:"size:255;not null;default:''"`
	AvatarURL    string    `gorm:"size:255;not null;default:'public/img/user/default.jpg'"`
	IsActive     bool      `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	OTPRequests   []OTPRequest
	Notifications []Notification
}

func SeedUsers() []User {
	now := time.Now()
	password, err := utils.HashString("admin12345")
	if err != nil {
		panic(err)
	}
	users := []User{
		{
			ID:           uuid.New(),
			Name:         "Super Admin",
			Email:        "super@admin.com",
			PasswordHash: password,
			Role:         "superadmin",
			Address:      "Kab. Gresik",
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    &now,
		},
		{
			ID:           uuid.New(),
			Name:         "Bayu Firmansyah",
			Email:        "bayu19fr@gmail.com",
			PasswordHash: password,
			Role:         "admin",
			Address:      "Kab. Gresik",
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    &now,
		},
		{
			ID:           uuid.New(),
			Name:         "Alif Dwi Rahman",
			Email:        "alifdwirahman.alf@gmail.com",
			PasswordHash: password,
			Role:         "admin",
			Address:      "Kota Jakarta",
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    &now,
		},
		{
			ID:           uuid.New(),
			Name:         "Fathoni Nur Habibi",
			Email:        "habibifathoni1509@gmail.com",
			PasswordHash: password,
			Role:         "admin",
			Address:      "Kota Jakarta",
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    &now,
		},
	}

	return users
}
