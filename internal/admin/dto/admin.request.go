package dto

type AttendancePeriodRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	By        int64  `json:"by" swaggerignore:"true"`
}

type PayrollRequest struct {
	ID       int64  `json:"id" swaggerignore:"true"`
	PeriodID int64  `json:"period_id"`
	Status   string `json:"status" swaggerignore:"true"`
	By       int64  `json:"by" swaggerignore:"true"`
}

type SummaryPayrollRequest struct {
	PayrollID int64 `json:"payroll_id"`
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
}
