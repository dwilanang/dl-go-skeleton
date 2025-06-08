package model

type AttendancePeriod struct {
	ID          int64  `json:"id"`
	StartDate   string `db:"start_date"`
	EndDate     string `db:"end_date"`
	IsProcessed bool   `db:"is_processed"`
	CreatedBy   int64  `db:"created_by"`
	CreatedAt   string `db:"created_at"`
	UpdatedBy   int64  `db:"updated_by"`
	UpdatedAt   string `db:"updated_at"`
}

type Payroll struct {
	ID          int64  `json:"id"`
	PeriodID    int64  `db:"period_id"`
	Status      string `db:"status"`
	ProcessedAt string `db:"processed_at"`
	CreatedBy   int64  `db:"created_by"`
	CreatedAt   string `db:"created_at"`
	UpdatedBy   int64  `db:"updated_by"`
	UpdatedAt   string `db:"updated_at"`
}

type PayrollItem struct {
	ID          int64  `db:"id"`
	Status      string `db:"status"`
	ProcessedAt string `db:"processed_at"`
	StartDate   string `db:"start_date"`
	EndDate     string `db:"end_date"`
	CreatedBy   int64  `db:"created_by"`
	CreatedAt   string `db:"created_at"`
	UpdatedBy   int64  `db:"updated_by"`
	UpdatedAt   string `db:"updated_at"`
}

type PeriodInfo struct {
	ID        int64  `db:"id"`
	StartDate string `db:"start_date"`
	EndDate   string `db:"end_date"`
}

type EmployeePayslip struct {
	ID             int64   `db:"id"`
	FullName       string  `db:"full_name"`
	BaseSalary     float64 `db:"base_salary"`
	AttendanceDays float64 `db:"attendance_days"`
	OvertimeHours  float64 `db:"overtime_hours"`
	OvertimePay    float64 `db:"overtime_pay"`
	Reimbursements float64 `db:"reimbursements"`
}
