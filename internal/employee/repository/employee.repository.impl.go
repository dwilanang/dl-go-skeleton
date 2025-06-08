package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dwilanang/psp/internal/employee/model"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateAttendance(a *model.Attendance) error {
	query := `
		INSERT INTO attendances
		(user_id, period_id, date_attendance, created_by, created_at, updated_by, updated_at) 
		VALUES ($1, $2, $3, $4, NOW(), $4, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRowx(
		query,
		a.UserID,
		a.PeriodID,
		a.DateAttendance,
		a.CreatedBy,
	).Scan(&a.ID, &a.CreatedAt)

	return err
}

func (r *repository) FetchAttendances(userID int64, limit, offset int) ([]*model.AttendanceItem, error) {
	query := `
		SELECT 
			a.id, 
			u.full_name,
			TO_CHAR(COALESCE(a.created_at, '0001-01-01 00:00:00'::timestamp), 'YYYY-MM-DD HH24:MI:SS') as attendances_at
		FROM 
			attendances a
		INNER JOIN
			users u ON a.user_id = u.id
		WHERE a.user_id = $1
		ORDER BY a.created_at DESC LIMIT $2 OFFSET $3
	`

	var attendances []*model.AttendanceItem
	err := r.db.Select(&attendances, query, userID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return attendances, nil
		}
		return attendances, err
	}

	return attendances, nil
}

func (r *repository) CountAttendances(userID int64) (int64, error) {
	var total int64
	err := r.db.Get(&total, `
		SELECT COUNT(*) FROM (
			SELECT id
			FROM attendances WHERE user_id = $1 GROUP BY id
		) sub
	`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return total, nil
		}
		return total, err
	}

	return total, nil
}

func (r *repository) HasSubmittedAttendance(userID int64, date time.Time) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM attendances 
			WHERE user_id = $1 AND date_attendance = $2
		)
	`
	err := r.db.Get(&exists, query, userID, date.Format("2006-01-02"))
	return exists, err
}

func (r *repository) CreateOvertimes(a *model.Overtime) error {
	query := `
		INSERT INTO overtimes
		(user_id, period_id, date_overtime, hours, created_by, created_at, updated_by, updated_at) 
		VALUES ($1, $2, $3, $4, $5, NOW(), $5, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRowx(
		query,
		a.UserID,
		a.PeriodID,
		a.DateOvertime,
		a.Hours,
		a.CreatedBy,
	).Scan(&a.ID, &a.CreatedAt)

	return err
}

func (r *repository) FetchOvertimes(userID int64, limit, offset int) ([]*model.OvertimeItem, error) {
	query := `
		SELECT 
			a.id, 
			u.full_name,
			a.hours,
			TO_CHAR(a.date_overtime, 'YYYY-MM-DD') as date_overtime
		FROM 
			overtimes a
		INNER JOIN
			users u ON a.user_id = u.id
		WHERE
			a.user_id = $1
		ORDER BY a.created_at DESC LIMIT $2 OFFSET $3
	`

	var overtimes []*model.OvertimeItem
	err := r.db.Select(&overtimes, query, userID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return overtimes, nil
		}
		return overtimes, err
	}

	return overtimes, nil
}

func (r *repository) CountOvertimes(userID int64) (int64, error) {
	var total int64
	err := r.db.Get(&total, `
		SELECT COUNT(*) FROM (
			SELECT id
			FROM overtimes WHERE user_id = $1 GROUP BY id
		) sub
	`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return total, nil
		}
		return total, err
	}

	return total, nil
}

func (r *repository) FetchReimbursement(userID int64, limit, offset int) ([]*model.ReimbursementItem, error) {
	query := `
		SELECT 
			a.id, 
			u.full_name,
			a.amount,
			a.description
		FROM 
			reimbursements a
		INNER JOIN
			users u ON a.user_id = u.id
		WHERE a.user_id = $1
		ORDER BY a.created_at DESC LIMIT $2 OFFSET $3
	`

	var reimbursements []*model.ReimbursementItem
	err := r.db.Select(&reimbursements, query, userID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reimbursements, nil
		}
		return reimbursements, err
	}

	return reimbursements, nil
}

func (r *repository) CountReimbursement(userID int64) (int64, error) {
	var total int64
	err := r.db.Get(&total, `
		SELECT COUNT(*) FROM (
			SELECT id
			FROM reimbursements WHERE user_id = $1 GROUP BY id
		) sub
	`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return total, nil
		}
		return total, err
	}

	return total, nil
}

func (r *repository) HasSubmittedOvertime(userID int64, date time.Time) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM overtimes 
			WHERE user_id = $1 AND date_overtime = $2
		)
	`
	err := r.db.Get(&exists, query, userID, date.Format("2006-01-02"))
	return exists, err
}

func (r *repository) CreateReimbursement(mr *model.Reimbursement) error {
	query := `
		INSERT INTO reimbursements
		(user_id, period_id, amount, description, created_by, created_at, updated_by, updated_at) 
		VALUES ($1, $2, $3, $4, $5, NOW(), $5, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRowx(
		query,
		mr.UserID,
		mr.PeriodID,
		mr.Amount,
		mr.Description,
		mr.CreatedBy,
	).Scan(&mr.ID, &mr.CreatedAt)

	return err
}

func (r *repository) GetPayrollByID(userID int64) (int64, error) {
	var periodID int64
	query := `
		SELECT period_id FROM payrolls WHERE id = $1 AND status = 'processed'
	`
	err := r.db.Get(&periodID, query, userID)
	if err != nil {
		return periodID, err
	}

	return periodID, nil
}

func (r *repository) GetEmployeeByID(userID int64) (model.Employee, error) {
	var employee model.Employee
	query := `
		SELECT e.full_name, us.amount AS base_salary FROM users e INNER JOIN user_salaries us ON us.user_id=e.id WHERE e.id = $1
	`
	err := r.db.Get(&employee, query, userID)
	if err != nil {
		return employee, err
	}

	return employee, nil
}

func (r *repository) CalculateAttendances(userID, periodID int64) (int64, error) {
	var attendanceDays int64
	query := `
		SELECT COUNT(DISTINCT id) FROM attendances WHERE user_id = $1 AND period_id = $2
	`
	err := r.db.Get(&attendanceDays, query, userID, periodID)
	if err != nil {
		return attendanceDays, err
	}

	return attendanceDays, nil
}

func (r *repository) CalculateOvertimes(userID, periodID int64) (float64, error) {
	var totalOvertimeHours float64
	query := `
		SELECT COALESCE(SUM(hours), 0) FROM overtimes WHERE user_id = $1 AND period_id = $2
	`
	err := r.db.Get(&totalOvertimeHours, query, userID, periodID)
	if err != nil {
		return totalOvertimeHours, err
	}

	return totalOvertimeHours, nil
}

func (r *repository) CalculateOvertimePay(userID, periodID int64, baseSalary float64) (float64, error) {
	var overtimePay float64
	query := `
		SELECT COALESCE(SUM(hours * ($3 / 160.0) * 2), 0) FROM overtimes WHERE user_id = $1 AND period_id = $2
	`
	err := r.db.Get(&overtimePay, query, userID, periodID, baseSalary)
	if err != nil {
		return overtimePay, err
	}

	return overtimePay, nil
}

func (r *repository) CalculateReimburse(userID, periodID int64) (float64, error) {
	var totalReimburse float64
	query := `
		SELECT COALESCE(SUM(amount), 0) FROM reimbursements WHERE user_id = $1 AND period_id = $2
	`
	err := r.db.Get(&totalReimburse, query, userID, periodID)
	if err != nil {
		return totalReimburse, err
	}

	return totalReimburse, nil
}
