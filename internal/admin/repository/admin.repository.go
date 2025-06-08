package repository

import (
	"time"

	"github.com/dwilanang/psp/internal/admin/model"
)

//go:generate mockgen -source=admin.repository.go -package=mocks -destination=mocks/mock_admin_repository.go

// Repository defines methods for managing attendance periods and payroll data storage and retrieval.
type Repository interface {
	// ValidateAttendancePeriodDate checks if the given attendance period dates (start and end) are valid.
	// Returns an error if the dates overlap with existing periods or are otherwise invalid.
	ValidateAttendancePeriodDate(startDate, endDate string) error

	// CreateAttendancePeriods inserts a new attendance period record into the data store.
	CreateAttendancePeriods(ap *model.AttendancePeriod) error

	// FetchAttendancePeriods retrieves a list of attendance period with pagination.
	FetchAttendancePeriods(limit, offset int) ([]*model.AttendancePeriod, error)

	// CountAttendancePeriods returns the total count of attendance period.
	CountAttendancePeriods() (int64, error)

	// CreatePayrolls inserts a new payroll record linked to an attendance period.
	CreatePayrolls(payroll *model.Payroll) error

	// FetchPayrolls retrieves a list of payrolls with pagination.
	FetchPayrolls(limit, offset int) ([]*model.PayrollItem, error)

	// CountPayrolls returns the total count of payroll.
	CountPayrolls() (int64, error)

	// UpdatePayrolls updates an existing payroll record with new data.
	UpdatePayrolls(payroll *model.Payroll) error

	// GetAttendancePeriodID retrieves the attendance period ID that contains the given date.
	// Returns the ID or an error if no period matches.
	GetAttendancePeriodID(date time.Time) (int64, error)

	// ValidatePayroll checks the status of payroll associated with the given attendance period ID.
	// Returns the status string or an error.
	ValidatePayroll(periodID int64) (string, error)

	// FindPayrollsByID retrieves the attendance period ID associated with a specific payroll ID.
	// Useful to link payrolls back to attendance periods.
	FindPayrollsByID(id int64) (int64, error)

	// FindAttendancePeriodsByID retrieves detailed period information for a given attendance period ID.
	// Returns a PeriodInfo struct or nil if not found.
	FindAttendancePeriodsByID(id int64) (*model.PeriodInfo, error)

	// CountEmployeeByPeriodID returns the total count of employees that have attendance or payroll data for a given period.
	CountEmployeeByPeriodID(periodID int64) (int64, error)

	// FindEmployeeByPeriodID retrieves a list of employee payroll details for a given attendance period ID with pagination.
	// Returns a slice of EmployeePayslip pointers or an error.
	FindEmployeeByPeriodID(periodID, limit, offset int64) ([]*model.EmployeePayslip, error)
}
