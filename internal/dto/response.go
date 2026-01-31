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

// EMAIL
type Email struct {
	Type string

	Username string
	Email    string
	Subject  string
	Body     string

	Code     string
	Password string
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

type UserProfile struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type ListAdmin struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Permissions map[string]bool `json:"permissions"`
}

// Auth
type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ValidateSession struct {
	SessionID uuid.UUID
	UserID    uuid.UUID
	Role      string
}

type CodeOTP struct {
	OTPToken  *uuid.UUID `json:"otp_token,omitempty"`
	Code      *string    `json:"code,omitempty"`
	ExpiredAt time.Time  `json:"expired_at"`
}

// Category
type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// global prod response
type ProductGlobalResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Category  *CategoryResponse        `json:"category,omitempty"`
	Inventory *InventoryGlobalResponse `json:"inventory,omitempty"`
}

// INVENTORY
type InventoryGlobalResponse struct {
	ID    int64  `json:"id"`
	Stock int    `json:"stock"`
	Unit  string `json:"unit"`
}

type InventoryResponse struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Stock     int       `json:"stock"`
	Unit      string    `json:"unit"`
	CreatedAt time.Time `json:"created_at"`

	Product ProductResponseInInventory `json:"product"`
}

type ProductResponseInInventory struct {
	ID          int64   `json:"id"`
	CategoryID  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
}
