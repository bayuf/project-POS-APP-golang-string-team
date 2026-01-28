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

type CreateCategory struct {
	Name string `json:"name" binding:"required"`
}

type CreateProduct struct {
	CategoryID  int64   `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	IsAvailable *bool   `json:"is_available,omitempty"`
}
