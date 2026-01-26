package dto

import "time"

// USER
type CreateUser struct {
	Name             string    `json:"name" binding:"required, min=3"`
	Email            string    `json:"email" binding:"required,email"`
	Phone            string    `json:"phone" binding:"required,numeric"`
	Role             string    `json:"role" binding:"required,oneof=admin staff"`
	Birthdate        time.Time `json:"birthdate" binding:"required"`
	Salary           int64     `json:"salary" binding:"required,min=0"`
	Address          string    `json:"address" binding:"required"`
	AdditionalDetail *string   `json:"additional_detail"`
	AvatarURL        *string   `json:"avatar_url"`
	ShiftStart       time.Time `json:"shift_start" binding:"required"`
	ShiftEnd         time.Time `json:"shift_end" binding:"required"`
}

type UpdateUser struct {
	Name             string  `json:"name" binding:"required, min=3"`
	Email            string  `json:"email" binding:"required,email"`
	Phone            string  `json:"phone" binding:"required,numeric"`
	Address          string  `json:"address" binding:"required"`
	AdditionalDetail *string `json:"additional_detail"`
}
