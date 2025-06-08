package model

import "time"

type User struct {
	ID           int64     `db:"id"`
	UUID         string    `db:"uuid"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	FullName     string    `db:"full_name"`
	RoleID       int64     `db:"role_id"`
	Role         string    `db:"role"`
	CreatedBy    int64     `db:"created_by"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedBy    int64     `db:"updated_by"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserSalary struct {
	ID            int64     `db:"id"`
	UserID        int64     `db:"user_id"`
	FullName      string    `db:"full_name"`
	Amount        float64   `db:"salary_amount"`
	EffectiveFrom string    `db:"effective_from"`
	CreatedBy     int64     `db:"created_by"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedBy     int64     `db:"updated_by"`
	UpdatedAt     time.Time `db:"updated_at"`
}
