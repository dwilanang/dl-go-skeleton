package model

type Attendance struct {
	ID             int64  `json:"id"`
	PeriodID       int64  `db:"period_id"`
	UserID         int64  `db:"user_id"`
	DateAttendance string `db:"date_attendance"`
	CreatedBy      int64  `db:"created_by"`
	CreatedAt      string `db:"created_at"`
	UpdatedBy      int64  `db:"updated_by"`
	UpdatedAt      string `db:"updated_at"`
}

type AttendanceItem struct {
	ID           int64  `json:"id"`
	FullName     string `db:"full_name"`
	AttendanceAt string `db:"attendances_at"`
}

type Overtime struct {
	ID           int64  `json:"id"`
	PeriodID     int64  `db:"period_id"`
	UserID       int64  `db:"user_id"`
	DateOvertime string `db:"date_overtime"`
	Hours        int    `db:"hours"`
	CreatedBy    int64  `db:"created_by"`
	CreatedAt    string `db:"created_at"`
	UpdatedBy    int64  `db:"updated_by"`
	UpdatedAt    string `db:"updated_at"`
}

type OvertimeItem struct {
	ID         int64   `json:"id"`
	FullName   string  `db:"full_name"`
	OvertimeAt string  `db:"date_overtime"`
	Hours      float64 `db:"hours"`
	CreatedAt  string  `db:"created_at"`
}

type Reimbursement struct {
	ID          int64  `json:"id"`
	PeriodID    int64  `db:"period_id"`
	UserID      int64  `db:"user_id"`
	Amount      int64  `db:"amount"`
	Description string `db:"description"`
	CreatedBy   int64  `db:"created_by"`
	CreatedAt   string `db:"created_at"`
	UpdatedBy   int64  `db:"updated_by"`
	UpdatedAt   string `db:"updated_at"`
}

type ReimbursementItem struct {
	ID          int64   `json:"id"`
	FullName    string  `db:"full_name"`
	Amount      float64 `db:"amount"`
	Description string  `db:"description"`
	CreatedAt   string  `db:"created_at"`
}

type Employee struct {
	FullName   string  `db:"full_name"`
	BaseSalary float64 `db:"base_salary"`
}
