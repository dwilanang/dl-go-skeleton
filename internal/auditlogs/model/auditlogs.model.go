package model

import "time"

type AuditLogs struct {
	ID        int64     `db:"id"`
	FullName  string    `db:"full_name"`
	Action    string    `db:"action"`
	TableName string    `db:"table_name"`
	RecordID  int64     `db:"record_id"`
	IpAddress string    `db:"ip_address"`
	RequestID string    `db:"request_id"`
	CreatedBy int       `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedBy int       `db:"updated_by"`
	UpdatedAt time.Time `db:"updated_at"`
}
