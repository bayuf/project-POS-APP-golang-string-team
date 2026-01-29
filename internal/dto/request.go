package dto

import (
	"github.com/google/uuid"
)

// USER
type CreateUser struct {
	Name             string  `json:"name" binding:"required"`
	Email            string  `json:"email" binding:"required,email"`
	Phone            string  `json:"phone" binding:"required,numeric"`
	Role             string  `json:"role" binding:"required,oneof=admin staff"`
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
	CategoryID  int64   `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	IsAvailable *bool   `json:"is_available,omitempty"`
}
