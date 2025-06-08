package service

import "github.com/dwilanang/psp/internal/auditlogs/dto"

type Service interface {
	GetAuditLog(page, limit int) (*dto.AuditLogResponse, error)
}
