package dto

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

// AUTH
type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}
