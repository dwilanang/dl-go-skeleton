package service

import "github.com/dwilanang/psp/internal/admin/dto"

//go:generate mockgen -source=admin.service.go -package=mocks -destination=mocks/mock_admin_service.go

// Service defines a set of operations related to managing attendance periods, payroll creation, processing, and payroll summaries.
type Service interface {
	// CreateAttendancePeriods creates a new attendance period based on the provided request.
	// Typically used at the start of a month to set the attendance timeframe.
	CreateAttendancePeriods(request *dto.AttendancePeriodRequest) error

	// GetAttendancePeriods List of attendance period
	GetAttendancePeriods(page, limit int) (*dto.AttendancePeriodResponse, error)

	// CreatePayrolls creates payroll entries for a specific period without running payroll calculations.
	// This usually records an empty payroll that can be processed later.
	CreatePayrolls(request *dto.PayrollRequest) error

	// GetPayrolls List of payroll
	GetPayrolls(page, limit int) (*dto.PayrollsResponse, error)

	// RunPayrolls executes payroll processing for a given period.
	// This includes calculations based on attendance, overtime, and reimbursements.
	RunPayrolls(request *dto.PayrollRequest) error

	// SummaryPayrolls returns a payroll summary containing employee payroll data for a specific period.
	// The summary includes base salary, attendance days, overtime, reimbursements, and take-home pay.
	SummaryPayrolls(request *dto.SummaryPayrollRequest) (*dto.PayrollSummaryResponse, error)
}
