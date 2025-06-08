package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dwilanang/psp/internal/admin/model"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) FetchAttendancePeriods(limit, offset int) ([]*model.AttendancePeriod, error) {
	query := `
		SELECT id, 
		TO_CHAR(start_date, 'YYYY-MM-DD') AS start_date, 
		TO_CHAR(end_date, 'YYYY-MM-DD')  AS end_date 
		FROM attendance_periods ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`

	var attendancePeriod []*model.AttendancePeriod
	err := r.db.Select(&attendancePeriod, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return attendancePeriod, nil
		}
		return attendancePeriod, err
	}

	return attendancePeriod, nil
}

func (r *repository) CountAttendancePeriods() (int64, error) {
	var total int64
	err := r.db.Get(&total, `
		SELECT COUNT(*) FROM (
			SELECT id
			FROM attendance_periods GROUP BY id
		) sub
	`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return total, nil
		}
		return total, err
	}

	return total, nil
}

func (r *repository) CreateAttendancePeriods(ap *model.AttendancePeriod) error {
	query := `
		INSERT INTO attendance_periods 
		(start_date, end_date, created_by, created_at, updated_by, updated_at) 
		VALUES ($1, $2, $3, NOW(), $3, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRowx(
		query,
		ap.StartDate,
		ap.EndDate,
		ap.CreatedBy,
	).Scan(&ap.ID, &ap.CreatedAt)

	return err
}

func (r *repository) CreatePayrolls(payroll *model.Payroll) error {
	query := `
		INSERT INTO payrolls 
		(period_id, status, created_by, created_at, updated_by, updated_at) 
		VALUES ($1, $2, $3, NOW(), $3,  NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRowx(
		query,
		payroll.PeriodID,
		payroll.Status,
		payroll.CreatedBy,
	).Scan(&payroll.ID, &payroll.CreatedAt)

	return err
}

func (r *repository) FetchPayrolls(limit, offset int) ([]*model.PayrollItem, error) {
	query := `
		SELECT 
			p.id, 
			ap.start_date,
			ap.end_date,
			TO_CHAR(ap.start_date, 'YYYY-MM-DD') AS start_date, 
			TO_CHAR(ap.end_date, 'YYYY-MM-DD')  AS end_date, 
			p.status,
			TO_CHAR(COALESCE(p.processed_at, '0001-01-01 00:00:00'::timestamp), 'YYYY-MM-DD HH24:MI:SS') as processed_at
		FROM 
			payrolls p
		INNER JOIN
			attendance_periods ap ON p.period_id = ap.id
		ORDER BY p.created_at DESC LIMIT $1 OFFSET $2
	`

	var payrolls []*model.PayrollItem
	err := r.db.Select(&payrolls, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return payrolls, nil
		}
		return payrolls, err
	}

	return payrolls, nil
}

func (r *repository) CountPayrolls() (int64, error) {
	var total int64
	err := r.db.Get(&total, `
		SELECT COUNT(*) FROM (
			SELECT id
			FROM payrolls GROUP BY id
		) sub
	`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return total, nil
		}
		return total, err
	}

	return total, nil
}

func (r *repository) UpdatePayrolls(payroll *model.Payroll) error {
	query := `
		UPDATE payrolls SET status = $1, processed_at = NOW(), updated_by = $2, updated_at = NOW() WHERE id = $3
		RETURNING id, created_at
	`

	err := r.db.QueryRowx(
		query,
		payroll.Status,
		payroll.UpdatedBy,
		payroll.ID,
	).Scan(&payroll.ID, &payroll.UpdatedAt)

	return err
}

func (r *repository) ValidateAttendancePeriodDate(startDate, endDate string) error {
	var count int
	r.db.Get(&count, `
		SELECT COUNT(1)
		FROM attendance_periods
		WHERE NOT (
			$2 < start_date OR $1 > end_date
		)
	`, startDate, endDate)

	if count > 0 {
		return errors.New("the period overlaps with another existing period")
	}

	return nil
}

func (r *repository) GetAttendancePeriodID(date time.Time) (int64, error) {
	var periodID int64

	err := r.db.Get(&periodID, `
		SELECT id
		FROM attendance_periods
		WHERE $1 BETWEEN start_date AND end_date
		LIMIT 1
	`, date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return periodID, nil
}

func (r *repository) ValidatePayroll(periodID int64) (string, error) {
	var status string
	err := r.db.Get(&status, `
		SELECT status
		FROM payrolls
		WHERE period_id = $1
	`, periodID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return status, nil
}

func (r *repository) FindAttendancePeriodsByID(id int64) (*model.PeriodInfo, error) {
	var mp model.PeriodInfo
	query := `
		SELECT 
			id, 
			TO_CHAR(start_date, 'YYYY-MM-DD') AS start_date, 
			TO_CHAR(end_date, 'YYYY-MM-DD')  AS end_date 
		FROM attendance_periods WHERE id = $1
	`
	err := r.db.Get(&mp, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &mp, nil
		}
		return &mp, err
	}

	return &mp, nil
}

func (r *repository) FindPayrollsByID(id int64) (int64, error) {
	var periodID int64
	err := r.db.Get(&periodID, `
		SELECT period_id FROM payrolls WHERE id = $1
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return periodID, nil
		}
		return periodID, err
	}

	return periodID, nil
}

func (r *repository) CountEmployeeByPeriodID(periodID int64) (int64, error) {
	var total int64
	err := r.db.Get(&total, `
		SELECT COUNT(*) FROM (
			SELECT e.id
			FROM users e
			LEFT JOIN attendances a ON a.user_id = e.id AND a.period_id = $1
			LEFT JOIN overtimes o ON o.user_id = e.id AND o.period_id = $1
			LEFT JOIN reimbursements r ON r.user_id = e.id AND r.period_id = $1
			GROUP BY e.id
			HAVING COUNT(a.id) > 0 OR COUNT(o.id) > 0 OR COUNT(r.id) > 0
		) subquery
	`, periodID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return total, nil
		}
		return total, err
	}

	return total, nil
}

func (r *repository) FindEmployeeByPeriodID(periodID, limit, offset int64) ([]*model.EmployeePayslip, error) {
	query := `
		WITH eligible_employees AS (
			SELECT e.id
			FROM users e
			LEFT JOIN attendances a ON a.user_id = e.id AND a.period_id = $1
			LEFT JOIN overtimes o ON o.user_id = e.id AND o.period_id = $1
			LEFT JOIN reimbursements r ON r.user_id = e.id AND r.period_id = $1
			GROUP BY e.id
			HAVING COUNT(a.id) > 0 OR COUNT(o.id) > 0 OR COUNT(r.id) > 0
		),
		paged_employees AS (
			SELECT id
			FROM eligible_employees
			ORDER BY id
			LIMIT $2 OFFSET $3
		)
		SELECT
			e.id,
			e.full_name,
			us.amount as base_salary,
			COUNT(DISTINCT a.id) AS attendance_days,
			COALESCE(SUM(o.hours), 0) AS overtime_hours,
			COALESCE(SUM(o.hours * (us.amount / 160.0) * 2), 0) AS overtime_pay,
			COALESCE(SUM(r.amount), 0) AS reimbursements
		FROM users e
		INNER JOIN user_salaries us ON us.user_id = e.id
		JOIN paged_employees pe ON pe.id = e.id
		LEFT JOIN attendances a ON a.user_id = e.id AND a.period_id = $1
		LEFT JOIN overtimes o ON o.user_id = e.id AND o.period_id = $1
		LEFT JOIN reimbursements r ON r.user_id = e.id AND r.period_id = $1
		GROUP BY e.id, e.full_name, us.amount
		ORDER BY e.id
	`
	var employeePayslip []*model.EmployeePayslip
	err := r.db.Select(&employeePayslip, query, periodID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return employeePayslip, nil
		}
		return employeePayslip, err
	}

	return employeePayslip, nil
}
