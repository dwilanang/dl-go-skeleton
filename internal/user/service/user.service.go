package service

import (
	"github.com/dwilanang/psp/internal/user/dto"
)

//go:generate mockgen -source=user.service.go -package=mocks -destination=mocks/mock_user_service.go

// Service defines the interface for business logic related to user management,
// including user registration and salary creation.
type Service interface {
	// Register handles the registration of a new user.
	// Param: request - a pointer to UserRequest DTO containing user registration data.
	// Returns a UserResponse DTO and an error if the registration fails.
	Register(request *dto.UserRequest) (dto.UserResponse, error)
}
