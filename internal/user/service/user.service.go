package service

import (
	"github.com/dwilanang/psp/internal/user/dto"
)

// Service defines the interface for business logic related to user management,
// including user registration and salary creation.
type Service interface {
	// Register handles the registration of a new user.
	// Param: request - a pointer to UserRequest DTO containing user registration data.
	// Returns a UserResponse DTO and an error if the registration fails.
	Register(request *dto.UserRequest) (dto.UserResponse, error)

	// CreateSalary creates a salary record for a registered user.
	// Param: request - a pointer to UserSalaryRequest DTO containing salary details.
	// Returns a UserSalaryResponse DTO and an error if the operation fails.
	CreateSalary(request *dto.UserSalaryRequest) (dto.UserSalaryResponse, error)
}
