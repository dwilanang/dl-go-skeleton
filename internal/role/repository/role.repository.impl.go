package repository

import (
	"fmt"

	"github.com/dwilanang/psp/internal/role/model"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) Fetch() ([]*model.Role, error) {
	var roles []*model.Role
	query := `
		SELECT 
			rs.id, rs.name, 
			COALESCE(rs.privilege, '') AS privilege,
			rs.created_by,
			us1.full_name as created_by_name,
			rs.updated_by,
			us2.full_name as updated_by_name,
			TO_CHAR(rs.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
			TO_CHAR(rs.created_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at
		FROM roles rs
		INNER JOIN users us1 ON(us1.id=rs.created_by)
		INNER JOIN users us2 ON(us2.id=rs.updated_by)
	`
	err := r.db.Select(&roles, query)
	if err != nil {
		fmt.Println("Fetch: ", err)
		return nil, err
	}
	return roles, nil
}

func (r *repository) FindByID(id int64) (*model.Role, error) {
	var role model.Role
	query := `
		SELECT 
			rs.id, rs.name, 
			COALESCE(rs.privilege, '') AS privilege,
			rs.created_by,
			us1.full_name as created_by_name,
			rs.updated_by,
			us2.full_name as updated_by_name,
			TO_CHAR(rs.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
			TO_CHAR(rs.created_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at
		FROM roles rs
		INNER JOIN users us1 ON(us1.id=rs.created_by)
		INNER JOIN users us2 ON(us2.id=rs.updated_by)
		WHERE rs.id = $1
	`
	err := r.db.Get(&role, query, id)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repository) Create(role *model.Role) error {
	query := `
		INSERT INTO roles (name, privilege, created_by, created_at, updated_by, updated_at) VALUES ($1, $2, $3, NOW(), $3, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRowx(
		query,
		role.Name,
		role.Privilege,
		role.CreatedBy,
	).Scan(&role.ID, &role.CreatedAt)

	return err
}

func (r *repository) Update(role *model.Role) error {
	query := `
		UPDATE roles SET name = $1, privilege = $2, updated_by = $3, updated_at = NOW() WHERE id = $4
		RETURNING id, created_at
	`
	return r.db.QueryRowx(
		query,
		role.Name,
		role.Privilege,
		role.UpdatedBy,
		role.ID,
	).Scan(&role.ID, &role.CreatedAt)
}

func (r *repository) Delete(id int64) error {
	query := `
		DELETE FROM roles WHERE id = $1
	`
	_, err := r.db.Exec(
		query,
		id,
	)

	return err
}
