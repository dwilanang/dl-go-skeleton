package role

import (
	"github.com/dwilanang/psp/internal/role/service"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	DBPostgres *sqlx.DB
	Service    service.Service
}
