package dto

import "time"

type Pagination struct {
	CurrentPage  int   `json:"current_page"`
	Limit        int   `json:"limit"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

// USER
type UserDetail struct {
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	BirthDate  time.Time  `json:"birth_date"`
	Role       string     `json:"role"`
	Address    string     `json:"address,omitempty"`
	Salary     int64      `json:"salary"`
	AvatarURL  string     `json:"avatar_url,omitempty"`
	ShiftStart *time.Time `json:"shift_start,omitempty"`
	ShiftEnd   *time.Time `json:"shift_end,omitempty"`
}
