package repository

import "github.com/dwilanang/psp/internal/auditlogs/model"

type Repository interface {
	FetchAuditLog(limit, offset int) ([]*model.AuditLogs, error)
	CountAuditLog() (int64, error)
}
