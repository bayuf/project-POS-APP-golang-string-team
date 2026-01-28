package dto

import (
	"time"

	"github.com/google/uuid"
)

type Pagination struct {
	CurrentPage  int   `json:"current_page"`
	Limit        int   `json:"limit"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

// USER
type UserDetail struct {
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	BirthDate  time.Time `json:"birth_date"`
	Role       string    `json:"role"`
	Address    string    `json:"address,omitempty"`
	Salary     int64     `json:"salary"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	ShiftStart string    `json:"shift_start,omitempty"`
	ShiftEnd   string    `json:"shift_end,omitempty"`
}

type UserLists struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Age       int       `json:"age"`
	Role      string    `json:"role"`
	Salary    int64     `json:"salary"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Timing    string    `json:"timing,omitempty"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Categorys struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Category *Category `json:"category,omitempty"`
}
