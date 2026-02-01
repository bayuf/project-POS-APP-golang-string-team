package dto

import (
	"github.com/google/uuid"
)

// USER
type CreateUser struct {
	Name             string `json:"name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Phone            string `json:"phone" binding:"required,numeric"`
	Role             string `json:"role" binding:"required,oneof=admin staff"`
	Permissions      map[string]bool
	Birthdate        string  `json:"birthdate" binding:"required"`
	Salary           int64   `json:"salary" binding:"required,min=0"`
	Address          string  `json:"address" binding:"required"`
	AdditionalDetail *string `json:"additional_detail"`
	AvatarURL        *string `json:"avatar_url"`
	ShiftStart       string  `json:"shift_start" binding:"required"`
	ShiftEnd         string  `json:"shift_end" binding:"required"`
}

type UpdateUser struct {
	Name             string  `json:"name" binding:"required"`
	Email            string  `json:"email" binding:"required,email"`
	Role             string  `json:"role" binding:"required,oneof=admin staff"`
	Phone            string  `json:"phone" binding:"required,numeric"`
	Salary           int64   `json:"salary" binding:"required,min=0"`
	Birthdate        string  `json:"birthdate" binding:"required"`
	ShiftStart       string  `json:"shift_start" binding:"required"`
	ShiftEnd         string  `json:"shift_end" binding:"required"`
	Address          string  `json:"address" binding:"required"`
	AdditionalDetail *string `json:"additional_detail"`
	AvatarURL        *string `json:"avatar_url"`
}

type UserFilterRequest struct {
	Page   int    `form:"page" binding:"min=1"`
	Limit  int    `form:"limit" binding:"min=1,max=100"`
	SortBy string `form:"sort_by"`
}

type UpdateUserProfile struct {
	Name            string `json:"name" binding:"omitempty"`
	Email           string `json:"email" binding:"omitempty,email"`
	Address         string `json:"address" binding:"omitempty"`
	NewPassword     string `json:"new_password" binding:"omitempty,min=5"`
	ConfirmPassword string `json:"confirm_password" binding:"omitempty,min=5"`
}

type UpdateUserPermissions struct {
	Permissions map[string]bool `json:"permissions" binding:"required"`
}

// AUTH
type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

type VerifyOTP struct {
	Email   string `json:"email" binding:"required,email"`
	OTPCode string `json:"otp_code" binding:"required"`
}

type UpdatePassword struct {
	Token           uuid.UUID `json:"token" binding:"required"`
	NewPassword     string    `json:"new_password" binding:"required,min=5"`
	ConfirmPassword string    `json:"confirm_password" binding:"required,min=5"`
}

// PRODUCT
type CreateCategory struct {
	Name string `json:"name" binding:"required"`
}

type CreateProduct struct {
	CategoryID  int64  `json:"category_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Price       int64  `json:"price" binding:"required,gte=0"`
	IsAvailable *bool  `json:"is_available,omitempty"`
}

// INVENTORY
type CreateInventory struct {
	ProductID int64  `json:"product_id" binding:"required"`
	Stock     int    `json:"stock" binding:"required,min=0"`
	Unit      string `json:"unit" binding:"omitempty,max=20"`
}

type UpdateInventory struct {
	Stock int    `json:"stock" binding:"required,min=0"`
	Unit  string `json:"unit" binding:"omitempty,max=20"`
}

type InventoryFilterRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

type SearchInventoryQuery struct {
	ProductName string `form:"product_name"`
	CategoryID  *int64 `form:"category_id"`
	MinStock    *int   `form:"min_stock"`
	MaxStock    *int   `form:"max_stock"`
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
}

// ORDERS
type OrderItem struct {
	ID       int64 `json:"item_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required,min=1"`
}

type Order struct {
	CustomerName string      `json:"customer_name" binding:"required"`
	TableID      int64       `json:"table_id" binding:"omitempty"`
	Orders       []OrderItem `json:"orders" binding:"required"`
// RESERVATION
type CreateReservationRequest struct {
	CustomerName    string `json:"customer_name" binding:"required,min=3,max=100"`
	TableID         int64  `json:"table_id" binding:"required"`
	ReservationTime string `json:"reservation_time" binding:"required"`
}

// type UpdateReservation struct {
// 	CustomerName    string `json:"customer_name" binding:"omitempty,min=3,max=100"`
// 	TableID         int64  `json:"table_id" binding:"omitempty"`
// 	ReservationTime string `json:"reservation_time" binding:"omitempty"`
// 	IsCancelled     *bool  `json:"is_cancelled" binding:"omitempty"`
// }

type UpdateReservationRequest struct {
	TableID     int64 `json:"table_id" binding:"required"`
	IsCancelled bool  `json:"is_cancelled"`
}

type ReservationFilterRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Status     string `form:"status" binding:"omitempty,oneof=active cancelled"`
	SearchName string `form:"search_name"`
	Date       string `form:"date"` // filter berdasarkan date tertentu
}
