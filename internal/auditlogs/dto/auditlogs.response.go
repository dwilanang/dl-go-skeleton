package dto

type AuditLogResponse struct {
	Data        []AuditLogData `json:"data"`
	TotalRecord int64          `json:"total_record"`
	TotalPages  int64          `json:"total_pages"`
	Limit       int            `json:"limit"`
	Page        int            `json:"page"`
}

type AuditLogData struct {
	ID        int64  `json:"id"`
	FullName  string `json:"full_name"`
	Action    string `json:"action"`
	Tablename string `json:"table_name"`
	RecordID  int64  `json:"record_id"`
	IpAddress string `json:"ip_address"`
	RequestID string `json:"request_id"`
}
