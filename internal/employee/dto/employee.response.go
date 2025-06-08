package dto

type EmployeeResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type EmployeePayslip struct {
	PayrollID      int64   `json:"id"`
	FullName       string  `json:"full_name"`
	BaseSalary     float64 `json:"base_salary"`
	AttendanceDays int64   `json:"attendance_days"`
	AttendancePay  float64 `json:"attendance_pay"`
	OvertimeHours  float64 `json:"overtime_hours"`
	OvertimePay    float64 `json:"overtime_pay"`
	Reimbursements float64 `json:"reimbursements"`
	TakeHomePay    float64 `json:"take_home_pay"`
}

type AttendanceResponse struct {
	Data        []AttendanceData `json:"data"`
	TotalRecord int64            `json:"total_record"`
	TotalPages  int64            `json:"total_pages"`
	Limit       int              `json:"limit"`
	Page        int              `json:"page"`
}

type AttendanceData struct {
	ID           int64  `json:"id"`
	FullName     string `json:"status"`
	AttendanceAt string `json:"date_attendance"`
}

type OvertimeResponse struct {
	Data        []OvertimeData `json:"data"`
	TotalRecord int64          `json:"total_record"`
	TotalPages  int64          `json:"total_pages"`
	Limit       int            `json:"limit"`
	Page        int            `json:"page"`
}

type OvertimeData struct {
	ID         int64  `json:"id"`
	FullName   string `json:"status"`
	OvertimeAt string `json:"overtime_at"`
	Hours      int    `json:"hourse"`
}

type ReimbursementResponse struct {
	Data        []ReimbursementData `json:"data"`
	TotalRecord int64               `json:"total_record"`
	TotalPages  int64               `json:"total_pages"`
	Limit       int                 `json:"limit"`
	Page        int                 `json:"page"`
}

type ReimbursementData struct {
	ID          int64   `json:"id"`
	FullName    string  `json:"full_name"`
	Amout       float64 `json:"amount"`
	Description string  `json:"description"`
}
