package dto

type AdminResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type PayrollSummaryResponse struct {
	PayrollID        int64             `json:"payroll_id"`
	Period           PeriodInfo        `json:"period"`
	Employees        []EmployeePayslip `json:"employees"`
	TotalTakeHomePay float64           `json:"total_take_home_pay"`
	TotalRecord      int64             `json:"total_record"`
	TotalPages       int64             `json:"total_pages"`
	CurrentPage      int               `json:"current_page"`
	Limit            int               `json:"limit"`
	Page             int64             `json:"page"`
}

type PeriodInfo struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type EmployeePayslip struct {
	ID             int64   `json:"id"`
	FullName       string  `json:"full_name"`
	BaseSalary     float64 `json:"base_salary"`
	AttendanceDays float64 `json:"attendance_days"`
	AttendancePay  float64 `json:"attendance_pay"`
	OvertimeHours  float64 `json:"overtime_hours"`
	OvertimePay    float64 `json:"overtime_pay"`
	Reimbursements float64 `json:"reimbursements"`
	TakeHomePay    float64 `json:"take_home_pay"`
}

type AttendancePeriodResponse struct {
	Data        []AttendancePeriodData `json:"data"`
	TotalRecord int64                  `json:"total_record"`
	TotalPages  int64                  `json:"total_pages"`
	Limit       int                    `json:"limit"`
	Page        int                    `json:"page"`
}

type AttendancePeriodData struct {
	ID        int64  `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type PayrollsResponse struct {
	Data        []PayrollsResponseData `json:"data"`
	TotalRecord int64                  `json:"total_record"`
	TotalPages  int64                  `json:"total_pages"`
	Limit       int                    `json:"limit"`
	Page        int                    `json:"page"`
}

type PayrollsResponseData struct {
	ID          int64  `json:"id"`
	Status      string `json:"status"`
	ProcessedAt string `json:"processed_at"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}
