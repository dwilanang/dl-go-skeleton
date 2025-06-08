package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dwilanang/psp/internal/admin/model"
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

func TestCreateAttendancePeriods(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	ap := &model.AttendancePeriod{
		StartDate: startDate,
		EndDate:   endDate,
	}

	rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now())

	mock.ExpectQuery("INSERT INTO attendance_periods").
		WithArgs(ap.StartDate, ap.EndDate, ap.CreatedBy).
		WillReturnRows(rows)

	err := repo.CreateAttendancePeriods(ap)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), ap.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePayrolls(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	now := time.Now()
	payroll := &model.Payroll{
		PeriodID:  10,
		Status:    "pending",
		CreatedBy: 1,
	}

	// Expect the insert query
	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO payrolls 
		(period_id, status, created_by, created_at, updated_by, updated_at) 
		VALUES ($1, $2, $3, NOW(), $3,  NOW())
		RETURNING id, created_at
	`)).
		WithArgs(payroll.PeriodID, payroll.Status, payroll.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow(2, now))

	err := repo.CreatePayrolls(payroll)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), payroll.ID)
	createdAtTime, err := time.Parse(time.RFC3339, payroll.CreatedAt)
	assert.NoError(t, err)
	assert.WithinDuration(t, now, createdAtTime, time.Second)

	// Ensure expectations met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdatePayrolls(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	now := time.Now()
	payroll := &model.Payroll{
		ID:        5,
		Status:    "PROCESSED",
		UpdatedBy: 1,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE payrolls SET status = $1, processed_at = NOW(), updated_by = $2, updated_at = NOW() WHERE id = $3
		RETURNING id, created_at
	`)).
		WithArgs(payroll.Status, payroll.UpdatedBy, payroll.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow(payroll.ID, now))

	err := repo.UpdatePayrolls(payroll)
	assert.NoError(t, err)
	assert.Equal(t, payroll.ID, payroll.ID)
	updatedAtTime, err := time.Parse(time.RFC3339, payroll.UpdatedAt)
	assert.NoError(t, err)
	assert.WithinDuration(t, now, updatedAtTime, time.Second)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestValidateAttendancePeriodDate(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	t.Run("no overlap", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(1\\)").
			WithArgs(startDate, endDate).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		err := repo.ValidateAttendancePeriodDate(startDate, endDate)
		assert.NoError(t, err)
	})

	t.Run("overlap found", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(1\\)").
			WithArgs(startDate, endDate).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		err := repo.ValidateAttendancePeriodDate(startDate, endDate)
		assert.EqualError(t, err, "the period overlaps with another existing period")
	})
}

func TestGetAttendancePeriodID(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	testDate := time.Now()

	t.Run("period found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id").
			WithArgs(testDate).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

		id, err := repo.GetAttendancePeriodID(testDate)
		assert.NoError(t, err)
		assert.Equal(t, int64(10), id)
	})

	t.Run("period not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id").
			WithArgs(testDate).
			WillReturnError(sql.ErrNoRows)

		id, err := repo.GetAttendancePeriodID(testDate)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), id)
	})
}

func TestValidatePayroll(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewRepository(db)

	periodID := int64(1)

	t.Run("payroll found", func(t *testing.T) {
		mock.ExpectQuery("SELECT status").
			WithArgs(periodID).
			WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("PENDING"))

		status, err := repo.ValidatePayroll(periodID)
		assert.NoError(t, err)
		assert.Equal(t, "PENDING", status)
	})

	t.Run("payroll not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT status").
			WithArgs(periodID).
			WillReturnError(sql.ErrNoRows)

		status, err := repo.ValidatePayroll(periodID)
		assert.NoError(t, err)
		assert.Equal(t, "", status)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT status").
			WithArgs(periodID).
			WillReturnError(errors.New("query error"))

		status, err := repo.ValidatePayroll(periodID)
		assert.EqualError(t, err, "query error")
		assert.Equal(t, "", status)
	})
}
