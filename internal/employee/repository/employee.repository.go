package repository

import (
	"time"

	"github.com/dwilanang/psp/internal/employee/model"
)

//go:generate mockgen -source=employee.repository.go -package=mocks -destination=mocks/mock_employee_repository.go

// Repository defines the data access layer contract for managing attendance, overtime, reimbursements,
// payroll, and employee-related queries in the database.
//
// Implementations of this interface are responsible for all direct interactions with the persistent storage.
type Repository interface {
	// CreateAttendance inserts a new attendance record into the database.
	//
	// Parameters:
	//   - a: pointer to an Attendance model containing the attendance details to be stored.
	//
	// Returns:
	//   - error: non-nil if the insertion fails.
	CreateAttendance(a *model.Attendance) error

	FetchAttendances(userID int64, limit, offset int) ([]*model.AttendanceItem, error)

	CountAttendances(userID int64) (int64, error)

	// HasSubmittedAttendance checks whether the specified user has already submitted attendance for a given date.
	//
	// Parameters:
	//   - userID: the unique identifier of the employee.
	//   - date: the date for which to check attendance submission.
	//
	// Returns:
	//   - bool: true if attendance has been submitted for that date; false otherwise.
	//   - error: non-nil if the check fails.
	HasSubmittedAttendance(userID int64, date time.Time) (bool, error)

	// CreateOvertimes inserts a new overtime record into the database.
	//
	// Parameters:
	//   - o: pointer to an Overtime model containing overtime details to be stored.
	//
	// Returns:
	//   - error: non-nil if the insertion fails.
	CreateOvertimes(o *model.Overtime) error

	FetchOvertimes(userID int64, limit, offset int) ([]*model.OvertimeItem, error)

	CountOvertimes(userID int64) (int64, error)

	// HasSubmittedOvertime checks whether the specified user has already submitted overtime for a given date.
	//
	// Parameters:
	//   - userID: the unique identifier of the employee.
	//   - date: the date for which to check overtime submission.
	//
	// Returns:
	//   - bool: true if overtime has been submitted for that date; false otherwise.
	//   - error: non-nil if the check fails.
	HasSubmittedOvertime(userID int64, date time.Time) (bool, error)

	// CreateReimbursement inserts a new reimbursement record into the database.
	//
	// Parameters:
	//   - r: pointer to a Reimbursement model containing reimbursement details to be stored.
	//
	// Returns:
	//   - error: non-nil if the insertion fails.
	CreateReimbursement(r *model.Reimbursement) error

	FetchReimbursement(userID int64, limit, offset int) ([]*model.ReimbursementItem, error)

	CountReimbursement(userID int64) (int64, error)

	// GetPayrollByID retrieves the payroll period ID associated with the given user.
	//
	// Parameters:
	//   - userID: the unique identifier of the employee.
	//
	// Returns:
	//   - int64: the payroll period ID.
	//   - error: non-nil if the retrieval fails.
	GetPayrollByID(userID int64) (int64, error)

	// GetEmployeeByID fetches the employee record corresponding to the given user ID.
	//
	// Parameters:
	//   - userID: the unique identifier of the employee.
	//
	// Returns:
	//   - model.Employee: the employee data model.
	//   - error: non-nil if the retrieval fails.
	GetEmployeeByID(userID int64) (model.Employee, error)

	// CalculateAttendances calculates the total attendance count for a user within a specified payroll period.
	//
	// Parameters:
	//   - userID: the employee's unique identifier.
	//   - periodID: the payroll period identifier.
	//
	// Returns:
	//   - int64: total attendance count.
	//   - error: non-nil if calculation fails.
	CalculateAttendances(userID, periodID int64) (int64, error)

	// CalculateOvertimes calculates the total overtime hours for a user within a specified payroll period.
	//
	// Parameters:
	//   - userID: the employee's unique identifier.
	//   - periodID: the payroll period identifier.
	//
	// Returns:
	//   - float64: total overtime hours.
	//   - error: non-nil if calculation fails.
	CalculateOvertimes(userID, periodID int64) (float64, error)

	// CalculateOvertimePay calculates the total overtime pay for a user within a specified payroll period,
	// based on the user's base salary.
	//
	// Parameters:
	//   - userID: the employee's unique identifier.
	//   - periodID: the payroll period identifier.
	//   - baseSalary: the base salary amount to use for pay calculation.
	//
	// Returns:
	//   - float64: total overtime pay amount.
	//   - error: non-nil if calculation fails.
	CalculateOvertimePay(userID, periodID int64, baseSalary float64) (float64, error)

	// CalculateReimburse calculates the total reimbursement amount for a user within a specified payroll period.
	//
	// Parameters:
	//   - userID: the employee's unique identifier.
	//   - periodID: the payroll period identifier.
	//
	// Returns:
	//   - float64: total reimbursement amount.
	//   - error: non-nil if calculation fails.
	CalculateReimburse(userID, periodID int64) (float64, error)
}
