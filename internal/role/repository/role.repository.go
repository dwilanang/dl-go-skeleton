package repository

import "github.com/dwilanang/psp/internal/role/model"

//go:generate mockgen -source=role.repository.go -package=mocks -destination=mocks/mock_role_repository.go

// Repository defines an interface for data operations related to the Role entity.
// It provides methods for basic CRUD operations and fetching by ID.
type Repository interface {
	// Fetch retrieves all available role records from the data store.
	// Returns a slice of Role pointers and an error if the operation fails.
	Fetch() ([]*model.Role, error)

	// Create inserts a new role record into the data store.
	// Param: role - a pointer to the Role entity to be created.
	// Returns an error if the operation fails.
	Create(role *model.Role) error

	// Update modifies an existing role record identified by the ID in the Role struct.
	// Param: role - a pointer to the Role entity with updated data.
	// Returns an error if the operation fails.
	Update(role *model.Role) error

	// Delete removes a role record from the data store by its ID.
	// Param: id - the ID of the Role to be deleted.
	// Returns an error if the operation fails.
	Delete(id int64) error

	// FindByID retrieves a role record by its ID.
	// Param: id - the ID of the Role to retrieve.
	// Returns a pointer to the Role and an error if the operation fails or the record is not found.
	FindByID(id int64) (*model.Role, error)
}
