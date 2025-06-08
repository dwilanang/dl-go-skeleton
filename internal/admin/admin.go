package admin

import (
	"github.com/dwilanang/psp/internal/admin/service"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	DBPostgres *sqlx.DB
	Service    service.Service
}
