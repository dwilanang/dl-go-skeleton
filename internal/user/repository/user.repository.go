package repository

import "github.com/dwilanang/psp/internal/user/model"

// Repository defines the interface for data access operations related to the User entity
// and associated salary records. It includes methods for retrieving users and creating users and salaries.
type Repository interface {
	// FindByUUID retrieves a user by their integer-based UUID.
	// Param: id - the UUID of the user.
	// Returns a pointer to the User model and an error if the user is not found or the query fails.
	FindByUUID(id int) (*model.User, error)

	// FindByUsername retrieves a user by their username.
	// Param: username - the username to search for.
	// Returns a pointer to the User model and an error if the user is not found or the query fails.
	FindByUsername(username string) (*model.User, error)

	// Create inserts a new user record into the data store.
	// Param: user - a pointer to the User model containing user data.
	// Returns an error if the insertion fails.
	Create(user *model.User) error

	// CreateSalary inserts a new salary record for a user into the data store.
	// Param: us - a pointer to the UserSalary model containing salary data.
	// Returns an error if the insertion fails.
	CreateSalary(us *model.UserSalary) error
}
