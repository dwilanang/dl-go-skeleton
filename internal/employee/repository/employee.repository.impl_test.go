package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dwilanang/psp/internal/employee/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "postgres")

	cleanup := func() {
		sqlxDB.Close()
	}

	return sqlxDB, mock, cleanup
}

func TestCreateAttendance(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	a := &model.Attendance{
		UserID:         1,
		PeriodID:       2,
		DateAttendance: time.Now().Format("2006-01-02 15:04:05"),
		CreatedBy:      100,
	}

	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO attendances
		(user_id, period_id, date_attendance, created_by, created_at, updated_by, updated_at) 
		VALUES ($1, $2, $3, $4, NOW(), $4, NOW())
		RETURNING id, created_at
	`)).WithArgs(a.UserID, a.PeriodID, a.DateAttendance, a.CreatedBy).
		WillReturnRows(rows)

	err := repo.CreateAttendance(a)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), a.ID)
}

func TestHasSubmittedAttendance(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	userID := int64(10)
	date := time.Now()
	expected := true

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT EXISTS (
			SELECT 1 FROM attendances 
			WHERE user_id = $1 AND date_attendance = $2
		)
	`)).WithArgs(userID, date.Format("2006-01-02")).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(expected))

	result, err := repo.HasSubmittedAttendance(userID, date)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestCreateOvertimes(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	o := &model.Overtime{
		UserID:       1,
		PeriodID:     2,
		DateOvertime: time.Now().Format("2006-01-02 15:04:05"),
		Hours:        3,
		CreatedBy:    100,
	}

	mock.ExpectQuery("INSERT INTO overtimes").
		WithArgs(o.UserID, o.PeriodID, o.DateOvertime, o.Hours, o.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	err := repo.CreateOvertimes(o)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), o.ID)
}

func TestHasSubmittedOvertime(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	userID := int64(10)
	date := time.Now()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(userID, date.Format("2006-01-02")).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	result, err := repo.HasSubmittedOvertime(userID, date)
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestCreateReimbursement(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	r := &model.Reimbursement{
		UserID:      1,
		PeriodID:    2,
		Amount:      150,
		Description: "Medical",
		CreatedBy:   100,
	}

	mock.ExpectQuery("INSERT INTO reimbursements").
		WithArgs(r.UserID, r.PeriodID, r.Amount, r.Description, r.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	err := repo.CreateReimbursement(r)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), r.ID)
}

func TestGetPayrollByID(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	userID := int64(5)
	expectedPeriodID := int64(20)

	mock.ExpectQuery("SELECT period_id FROM payrolls").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"period_id"}).AddRow(expectedPeriodID))

	id, err := repo.GetPayrollByID(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPeriodID, id)
}

func TestGetEmployeeByID(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	userID := int64(1)
	fullName := "John Doe"
	baseSalary := 1000000.0

	mock.ExpectQuery("SELECT e.full_name, us.amount").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"full_name", "base_salary"}).AddRow(fullName, baseSalary))

	emp, err := repo.GetEmployeeByID(userID)
	assert.NoError(t, err)
	assert.Equal(t, fullName, emp.FullName)
	assert.Equal(t, baseSalary, emp.BaseSalary)
}

func TestCalculateAttendances(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	userID := int64(1)
	periodID := int64(2)
	expected := int64(15)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(userID, periodID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expected))

	count, err := repo.CalculateAttendances(userID, periodID)
	assert.NoError(t, err)
	assert.Equal(t, expected, count)
}

func TestCalculateOvertimes(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	userID := int64(1)
	periodID := int64(2)
	expected := 12.5

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(hours\\)").
		WithArgs(userID, periodID).
		WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(expected))

	sum, err := repo.CalculateOvertimes(userID, periodID)
	assert.NoError(t, err)
	assert.Equal(t, expected, sum)
}

func TestCalculateOvertimePay(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	userID := int64(1)
	periodID := int64(2)
	baseSalary := 1600000.0
	expectedPay := 400000.0

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(hours \\*").
		WithArgs(userID, periodID, baseSalary).
		WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(expectedPay))

	pay, err := repo.CalculateOvertimePay(userID, periodID, baseSalary)
	assert.NoError(t, err)
	assert.Equal(t, expectedPay, pay)
}

func TestCalculateReimburse(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	userID := int64(1)
	periodID := int64(2)
	expected := 200000.0

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(amount\\)").
		WithArgs(userID, periodID).
		WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(expected))

	sum, err := repo.CalculateReimburse(userID, periodID)
	assert.NoError(t, err)
	assert.Equal(t, expected, sum)
}
