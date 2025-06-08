package service

import "github.com/dwilanang/psp/internal/employee/dto"

//go:generate mockgen -source=employee.service.go -package=mocks -destination=mocks/mock_employee_service.go

// Service defines the business logic contract for managing employee attendance, overtime, reimbursements,
// and payslip generation within the application.
//
// Implementations of this interface are responsible for handling all operations related to
// employee payroll processing and related records.
type Service interface {
	// CreateAttendance processes and stores a new attendance record based on the provided request data.
	//
	// Parameters:
	//   - request: a pointer to AttendanceRequest containing attendance details to be saved.
	//
	// Returns:
	//   - error: non-nil if the operation failed.
	CreateAttendance(request *dto.AttendanceRequest) error

	GetAttendance(userID int64, page, limit int) (*dto.AttendanceResponse, error)

	// CreateOvertime processes and stores a new overtime record based on the provided request data.
	//
	// Parameters:
	//   - request: a pointer to OvertimeRequest containing overtime details to be saved.
	//
	// Returns:
	//   - error: non-nil if the operation failed.
	CreateOvertime(request *dto.OvertimeRequest) error

	GetOvertime(userID int64, page, limit int) (*dto.OvertimeResponse, error)

	// CreateReimbursement processes and stores a new reimbursement record based on the provided request data.
	//
	// Parameters:
	//   - request: a pointer to ReimbursementRequest containing reimbursement details to be saved.
	//
	// Returns:
	//   - error: non-nil if the operation failed.
	CreateReimbursement(request *dto.ReimbursementRequest) error

	GetReimbursement(userID int64, page, limit int) (*dto.ReimbursementResponse, error)

	// GenerateEmployeePayslip generates and returns the payslip data for a given employee and payroll period.
	//
	// Parameters:
	//   - userID: the unique identifier of the employee.
	//   - payrollID: the unique identifier of the payroll period.
	//
	// Returns:
	//   - *EmployeePayslip: pointer to the generated payslip DTO.
	//   - error: non-nil if the generation failed.
	GenerateEmployeePayslip(userID, payrollID int64) (*dto.EmployeePayslip, error)
}
