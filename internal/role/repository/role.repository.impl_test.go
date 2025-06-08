package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dwilanang/psp/internal/role/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupDBMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	return sqlx.NewDb(db, "postgres"), mock
}

func TestRepository_Fetch(t *testing.T) {
	db, mock := setupDBMock(t)
	defer db.Close()

	query := regexp.QuoteMeta(`
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
	`)

	rows := sqlmock.NewRows([]string{
		"id", "name", "privilege", "created_by", "created_by_name",
		"updated_by", "updated_by_name", "created_at", "updated_at",
	}).AddRow(1, "Admin", "all", 1, "Super Admin", 1, "Super Admin", "2024-01-01 00:00:00", "2024-01-01 00:00:00")

	mock.ExpectQuery(query).WillReturnRows(rows)

	repo := NewRepository(db)
	result, err := repo.Fetch()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Admin", result[0].Name)
}

func TestRepository_FindByID(t *testing.T) {
	db, mock := setupDBMock(t)
	defer db.Close()

	query := regexp.QuoteMeta(`
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
	`)

	rows := sqlmock.NewRows([]string{
		"id", "name", "privilege", "created_by", "created_by_name",
		"updated_by", "updated_by_name", "created_at", "updated_at",
	}).AddRow(1, "Admin", "all", 1, "Super Admin", 1, "Super Admin", "2024-01-01 00:00:00", "2024-01-01 00:00:00")

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	repo := NewRepository(db)
	result, err := repo.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "Admin", result.Name)
}

func TestRepository_Create(t *testing.T) {
	db, mock := setupDBMock(t)
	defer db.Close()

	query := regexp.QuoteMeta(`
		INSERT INTO roles (name, privilege, created_by, created_at, updated_by, updated_at) VALUES ($1, $2, $3, NOW(), $3, NOW())
		RETURNING id, created_at
	`)

	createdAt := time.Now()
	mock.ExpectQuery(query).
		WithArgs("Admin", "all", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, createdAt))

	repo := NewRepository(db)
	role := &model.Role{
		Name:      "Admin",
		Privilege: "all",
		CreatedBy: 1,
	}

	err := repo.Create(role)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), role.ID)
}

func TestRepository_Update(t *testing.T) {
	db, mock := setupDBMock(t)
	defer db.Close()

	query := regexp.QuoteMeta(`
		UPDATE roles SET name = $1, privilege = $2, updated_by = $3, updated_at = NOW() WHERE id = $4
		RETURNING id, created_at
	`)

	createdAt := time.Now()
	mock.ExpectQuery(query).
		WithArgs("Updated", "write", int64(2), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, createdAt))

	repo := NewRepository(db)
	role := &model.Role{
		ID:        1,
		Name:      "Updated",
		Privilege: "write",
		UpdatedBy: 2,
	}

	err := repo.Update(role)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), role.ID)
}

func TestRepository_Delete(t *testing.T) {
	db, mock := setupDBMock(t)
	defer db.Close()

	query := regexp.QuoteMeta(`DELETE FROM roles WHERE id = $1`)

	mock.ExpectExec(query).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewRepository(db)
	err := repo.Delete(1)
	assert.NoError(t, err)
}
