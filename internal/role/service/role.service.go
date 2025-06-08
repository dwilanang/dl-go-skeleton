package service

import (
	"github.com/dwilanang/psp/internal/role/dto"
)

//go:generate mockgen -source=role.service.go -package=rolemock -destination=../../rolemock/mock_role_service.go

// Service defines the interface for business logic related to the Role entity.
// It provides methods to retrieve, create, update, and delete roles.
type Service interface {
	// GetAll retrieves all roles and returns them in a structured response format.
	// Returns a RoleResponse DTO and an error if the operation fails.
	GetAll() (dto.RoleResponse, error)

	// Create adds a new role based on the provided request data.
	// Param: request - a pointer to RoleRequest DTO containing role details.
	// Returns an error if the creation fails.
	Create(request *dto.RoleRequest) error

	// Update modifies an existing role using the provided request data.
	// Param: request - a pointer to RoleRequest DTO with updated role information.
	// Returns an error if the update fails.
	Update(request *dto.RoleRequest) error

	// Delete removes a role identified by the given ID.
	// Param: id - the ID of the role to be deleted.
	// Returns an error if the deletion fails.
	Delete(id int64) error
}
