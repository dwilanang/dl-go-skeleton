package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/dwilanang/psp/internal/auth/model"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AuditMiddleware struct {
	DB *sqlx.DB
}

func NewAuditMiddleware(db *sqlx.DB) *AuditMiddleware {
	return &AuditMiddleware{DB: db}
}

func (a *AuditMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simpan request body jika perlu
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // reassign ke request
		}

		// Lanjutkan ke handler
		c.Next()

		// Hanya audit POST/PUT/DELETE
		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodPut && c.Request.Method != http.MethodDelete {
			return
		}

		// Ambil user dari context
		user, exists := c.Get("user")
		if !exists {
			return
		}

		claims, ok := user.(*model.TokenClaims)
		if !ok {
			return
		}

		// Ambil informasi
		action := c.Request.Method
		ip := c.ClientIP()
		requestID := c.Writer.Header().Get("X-Request-ID") // bisa dari header atau generate di middleware lain
		userID := claims.ID

		// Ambil nama resource dari route
		table := parseTableName(c.FullPath()) // misalnya: /api/v1/attendances/:id -> "attendances"

		// Ambil ID dari parameter kalau ada
		recordID := parseIDParam(c)

		// Simpan log ke database
		query := `
			INSERT INTO audit_logs (user_id, action, table_name, record_id, ip_address, request_id, created_by, updated_by, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $1, $1, NOW(), NOW())
		`

		_, err := a.DB.Exec(query, userID, action, table, recordID, ip, requestID)
		if err != nil {
			// log error, jangan return ke user
			c.Error(err)
		}
	}
}

func parseTableName(path string) string {
	// Sederhana: ambil segmen terakhir plural, contoh /api/v1/attendances/:id -> "attendances"
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && !strings.HasPrefix(parts[i], ":") {
			return parts[i]
		}
	}
	return "unknown"
}

func parseIDParam(c *gin.Context) int64 {
	idStr := c.Param("id")
	if idStr == "" {
		return 0
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0
	}
	return id
}
