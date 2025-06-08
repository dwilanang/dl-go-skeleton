package registry

import (
	"github.com/dwilanang/psp/config"
	"github.com/dwilanang/psp/internal/admin"
	adminhandler "github.com/dwilanang/psp/internal/admin/handler"
	adminrepository "github.com/dwilanang/psp/internal/admin/repository"
	adminservice "github.com/dwilanang/psp/internal/admin/service"
	"github.com/dwilanang/psp/internal/auditlogs"
	authhandler "github.com/dwilanang/psp/internal/auth/handler"
	authservice "github.com/dwilanang/psp/internal/auth/service"
	"github.com/dwilanang/psp/internal/role"
	rolehandler "github.com/dwilanang/psp/internal/role/handler"
	rolerepository "github.com/dwilanang/psp/internal/role/repository"
	roleservice "github.com/dwilanang/psp/internal/role/service"
	"github.com/dwilanang/psp/internal/user"
	userhandler "github.com/dwilanang/psp/internal/user/handler"
	userrepository "github.com/dwilanang/psp/internal/user/repository"
	userservice "github.com/dwilanang/psp/internal/user/service"

	"github.com/dwilanang/psp/internal/employee"
	employeehandler "github.com/dwilanang/psp/internal/employee/handler"
	employeerepository "github.com/dwilanang/psp/internal/employee/repository"
	employeeservice "github.com/dwilanang/psp/internal/employee/service"

	auditlogshandler "github.com/dwilanang/psp/internal/auditlogs/handler"
	auditlogsrepository "github.com/dwilanang/psp/internal/auditlogs/repository"
	auditlogsservice "github.com/dwilanang/psp/internal/auditlogs/service"

	"github.com/jmoiron/sqlx"
)

// Registry acts as a centralized dependency container.
// It provides factory methods to construct application components (handlers/services/repos)
// with their required dependencies injected, such as database connections and configuration settings.
type Registry struct {
	db  *sqlx.DB       // db represents a PostgreSQL database connection managed via sqlx.
	cfg *config.Config // cfg holds global configuration values used across services (e.g., JWT secret, environment settings).
}

// NewRegistry creates a new instance of the Registry.
// This function should be called once during application startup, providing it with
// the configuration and database connection to be used throughout the application.
func NewRegistry(cfg *config.Config, db *sqlx.DB) *Registry {
	return &Registry{db: db, cfg: cfg}
}

// NewAuthHandler returns a fully-initialized AuthHandler.
// It sets up the user repository and authentication service (which may handle login, JWT generation, etc.),
// and injects them into the auth handler.
func (r *Registry) NewAuthHandler() *authhandler.Handler {
	repo := userrepository.NewRepository(r.db)     // Repository to interact with user-related DB operations
	authSvc := authservice.NewService(r.cfg, repo) // Service encapsulating authentication logic
	return authhandler.NewHandler(authSvc)         // HTTP handler for auth-related routes
}

// NewRoleHandler returns a fully-initialized RoleHandler.
// It builds the role repository and service, and constructs the handler to manage role-related routes.
func (r *Registry) NewRoleHandler() *rolehandler.Handler {
	repo := rolerepository.NewRepository(r.db) // Repository to manage roles in the database
	svc := roleservice.NewService(repo)        // Business logic for managing roles
	return rolehandler.NewHandler(role.Dependencies{
		DBPostgres: r.db,
		Service:    svc,
	})
}

// NewUserHandler returns a fully-initialized UserHandler.
// This handler manages routes related to user management such as listing, creating, or updating users.
func (r *Registry) NewUserHandler() *userhandler.Handler {
	repo := userrepository.NewRepository(r.db) // Reusable user repository
	svc := userservice.NewService(repo)        // Service to handle user-related business logic
	return userhandler.NewHandler(user.Dependencies{
		DBPostgres: r.db,
		Service:    svc,
	})
}

// NewAdminHandler returns a fully-initialized AdminHandler.
// It provides access to admin-specific features such as admin account management.
func (r *Registry) NewAdminHandler() *adminhandler.Handler {
	repo := adminrepository.NewRepository(r.db) // Repository for admin data operations
	svc := adminservice.NewService(repo)        // Service layer encapsulating admin logic
	return adminhandler.NewHandler(admin.Dependencies{
		DBPostgres: r.db,
		Service:    svc,
	})
}

// NewEmployeeHandler returns a fully-initialized EmployeeHandler.
// This handler manages employee-related features, potentially requiring access to both employee and admin data.
func (r *Registry) NewEmployeeHandler() *employeehandler.Handler {
	adminRepo := adminrepository.NewRepository(r.db)           // Admin repository needed for employee context
	employeeRepo := employeerepository.NewRepository(r.db)     // Repository for employee-specific data
	svc := employeeservice.NewService(employeeRepo, adminRepo) // Business logic layer combining employee and admin data
	return employeehandler.NewHandler(employee.Dependencies{
		DBPostgres: r.db,
		Service:    svc,
	})
}

func (r *Registry) NewAuditLogsHandler() *auditlogshandler.Handler {
	auditlogsRepo := auditlogsrepository.NewRepository(r.db)
	svc := auditlogsservice.NewService(auditlogsRepo)
	return auditlogshandler.NewHandler(auditlogs.Dependencies{
		DBPostgres: r.db,
		Service:    svc,
	})
}
