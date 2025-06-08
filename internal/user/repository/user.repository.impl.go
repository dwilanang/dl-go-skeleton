package repository

import (
	"database/sql"
	"errors"

	"github.com/dwilanang/psp/internal/user/model"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) FindByUUID(id int) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE uuid = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	query := `
		SELECT 
			u.id,
			u.uuid,
			u.password_hash,
			r.name AS role
		FROM users u
		INNER JOIN roles r ON(u.role_id=r.id)
		WHERE u.username = $1
	`
	err := r.db.Get(&user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("failed")
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) Create(user *model.User) error {
	query := `
		INSERT INTO users (uuid, username, password_hash, full_name, role_id, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRowx(
		query,
		user.UUID,
		user.Username,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.CreatedBy,
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *repository) CreateSalary(us *model.UserSalary) error {
	query := `
		INSERT INTO user_salaries (user_id, amount, effective_from, created_by, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRowx(
		query,
		us.UserID,
		us.Amount,
		us.EffectiveFrom,
		us.CreatedBy,
	).Scan(&us.ID, &us.CreatedAt)
}
