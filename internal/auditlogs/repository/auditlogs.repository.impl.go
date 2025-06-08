package repository

import (
	"database/sql"
	"errors"

	"github.com/dwilanang/psp/internal/auditlogs/model"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) FetchAuditLog(limit, offset int) ([]*model.AuditLogs, error) {
	query := `
		SELECT 
			al.id,
			al.action,
			al.table_name,
			al.record_id,
			al.request_id,
			al.ip_address,
			u.full_name
		FROM 
			audit_logs al
		INNER JOIN
			users u ON al.user_id = u.id
		ORDER BY al.created_at DESC LIMIT $1 OFFSET $2
	`
	var results []*model.AuditLogs
	err := r.db.Select(&results, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return results, nil
		}
		return results, err
	}
	return results, nil
}

func (r *repository) CountAuditLog() (int64, error) {
	var total int64
	err := r.db.Get(&total, `
		SELECT COUNT(*) FROM (
			SELECT id
			FROM audit_logs GROUP BY id
		) sub
	`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return total, nil
		}
		return total, err
	}

	return total, nil
}
