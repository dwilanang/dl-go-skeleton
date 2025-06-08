package repository

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dwilanang/psp/internal/user/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	return sqlx.NewDb(db, "postgres"), mock, func() { db.Close() }
}

func TestFindByUUID_Success(t *testing.T) {
	db, mock, close := setupMockDB(t)
	defer close()

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "uuid", "username", "password_hash", "full_name", "role_id"}).
		AddRow(1, "uuid-123", "johndoe", "hashedpass", "John Doe", 2)

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.FindByUUID(1)
	fmt.Println(err)
	assert.NoError(t, err)
	assert.Equal(t, "johndoe", user.Username)
}

func TestFindByUsername_Success(t *testing.T) {
	db, mock, close := setupMockDB(t)
	defer close()

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "uuid", "password_hash", "role"}).
		AddRow(1, "uuid-123", "hashedpass", "admin")

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT 
			u.id,
			u.uuid,
			u.password_hash,
			r.name AS role
		FROM users u
		INNER JOIN roles r ON(u.role_id=r.id)
		WHERE u.username = $1
	`)).
		WithArgs("johndoe").
		WillReturnRows(rows)

	user, err := repo.FindByUsername("johndoe")

	assert.NoError(t, err)
	assert.Equal(t, "uuid-123", user.UUID)
	assert.Equal(t, "admin", user.Role)
}

func TestCreate_Success(t *testing.T) {
	db, mock, close := setupMockDB(t)
	defer close()

	repo := NewRepository(db)

	user := &model.User{
		UUID:         "uuid-123",
		Username:     "johndoe",
		PasswordHash: "hashedpass",
		FullName:     "John Doe",
		RoleID:       2,
		CreatedBy:    1,
	}

	createdAt := time.Now()
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.UUID, user.Username, user.PasswordHash, user.FullName, user.RoleID, user.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(10, createdAt))

	err := repo.Create(user)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), user.ID)
	assert.WithinDuration(t, createdAt, user.CreatedAt, time.Second)
}

func TestCreateSalary_Success(t *testing.T) {
	db, mock, close := setupMockDB(t)
	defer close()

	repo := NewRepository(db)

	salary := &model.UserSalary{
		UserID:        1,
		Amount:        5000000,
		EffectiveFrom: time.Now().Format("2006-01-02 15:04:05"),
		CreatedBy:     1,
	}

	createdAt := time.Now()
	mock.ExpectQuery("INSERT INTO user_salaries").
		WithArgs(salary.UserID, salary.Amount, salary.EffectiveFrom, salary.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(101, createdAt))

	err := repo.CreateSalary(salary)

	assert.NoError(t, err)
	assert.Equal(t, int64(101), salary.ID)
	assert.WithinDuration(t, createdAt, salary.CreatedAt, time.Second)
}
