package model

type Role struct {
	ID            int64  `json:"id"`
	Name          string `db:"name"`
	Privilege     string `db:"privilege"`
	CreatedBy     int64  `db:"created_by"`
	CreatedByName string `db:"created_by_name"`
	CreatedAt     string `db:"created_at"`
	UpdatedBy     int64  `db:"updated_by"`
	UpdatedByName string `db:"updated_by_name"`
	UpdatedAt     string `db:"updated_at"`
}
