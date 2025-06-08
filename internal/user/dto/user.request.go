package dto

type UserRequest struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	RoleID   int64  `json:"role_id" binding:"required"`
	By       int64  `json:"by"`
}

type UserSalaryRequest struct {
	UserID        int64   `json:"user_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	EffectiveFrom string  `json:"effective_from" binding:"required"`
	By            int64   `json:"by" swaggerignore:"true"`
}
