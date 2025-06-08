package auditlogs

import (
	"github.com/dwilanang/psp/internal/auditlogs/service"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	DBPostgres *sqlx.DB
	Service    service.Service
}
