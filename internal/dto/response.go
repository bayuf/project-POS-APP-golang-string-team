package dto

type Pagination struct {
	CurrentPage  int   `json:"current_page"`
	Limit        int   `json:"limit"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

// USER
type UserDetail struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Address   string `json:"address,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}
