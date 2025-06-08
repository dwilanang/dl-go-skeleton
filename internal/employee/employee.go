package employee

import (
	"github.com/dwilanang/psp/internal/employee/service"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	DBPostgres *sqlx.DB
	Service    service.Service
}
