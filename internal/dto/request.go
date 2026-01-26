package dto

// USER
type CreateUser struct {
	Name     string `json:"name" binding:"required, min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=admin staff"`
	Address  string `json:"address" binding:"required"`
}
