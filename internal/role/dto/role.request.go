package dto

type RoleRequest struct {
	ID        int64  `json:"id" swaggerignore:"true"`
	Name      string `json:"name" binding:"required"`
	Privilege string `json:"privilege" binding:"required"`
	By        int64  `json:"by" swaggerignore:"true"`
}
