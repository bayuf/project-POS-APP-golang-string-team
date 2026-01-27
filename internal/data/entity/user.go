package entity

import (
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Name         string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:100;not null;unique"`
	PasswordHash string    `gorm:"type:text;not null;default:''"`
	Phone        string    `gorm:"size:20;not null;default:''"`
	Role         string    `gorm:"size:20;not null"` // superadmin, admin, staff
	AvatarURL    string    `gorm:"size:255;not null;default:'public/img/user/default.jpg'"`
	BirthDate    time.Time `gorm:"type:date"`
	Salary       int64     `gorm:"not null;default:0"`

	ShiftStart string `gorm:"type:time"`
	ShiftEnd   string `gorm:"type:time"`

	Address          string  `gorm:"text;not null;default:''"`
	AdditionalDetail *string `gorm:"text;"`

	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time      `gorm:"type:timestamptz"`
	UpdatedAt time.Time      `gorm:"type:timestamptz"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index"`

	OTPRequests   []OTPRequest
	Notifications []Notification
	Session       []Session `gorm:"foreignKey:UserID"`
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
			Phone:        "089111111111",
			Role:         "superadmin",
			Address:      "Kab. Gresik",
			BirthDate:    time.Date(2000, 5, 19, 0, 0, 0, 0, time.UTC),
			Salary:       10000000,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New(),
			Name:         "Bayu Firmansyah",
			Email:        "bayu19fr@gmail.com",
			PasswordHash: password,
			Phone:        "089111111111",
			Role:         "admin",
			Address:      "Kab. Gresik",
			BirthDate:    time.Date(2000, 5, 19, 0, 0, 0, 0, time.UTC),
			Salary:       10000000,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New(),
			Name:         "Alif Dwi Rahman",
			Email:        "alifdwirahman.alf@gmail.com",
			PasswordHash: password,
			Phone:        "089111111111",
			Role:         "admin",
			Address:      "Kota Jakarta",
			BirthDate:    time.Date(2000, 5, 19, 0, 0, 0, 0, time.UTC),
			Salary:       10000000,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New(),
			Name:         "Fathoni Nur Habibi",
			Email:        "habibifathoni1509@gmail.com",
			PasswordHash: password,
			Phone:        "089111111111",
			Role:         "admin",
			Address:      "Kota Jakarta",
			BirthDate:    time.Date(2000, 5, 19, 0, 0, 0, 0, time.UTC),
			Salary:       10000000,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	return users
}
