package user

import (
	"github.com/dwilanang/psp/internal/user/service"
	"github.com/jmoiron/sqlx"
)

type Dependencies struct {
	DBPostgres *sqlx.DB
	Service    service.Service
}
