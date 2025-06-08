package logger

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	ContextKeyRequestID = "request_id"
	ContextKeyIPAddress = "ip_address"
)

// RequestLogger middleware untuk log + inject request_id dan ip
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Ambil atau buat request ID
		reqID := c.Request.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// Ambil IP
		ip := getClientIP(c)

		// Simpan ke context
		c.Set(ContextKeyRequestID, reqID)
		c.Set(ContextKeyIPAddress, ip)

		// Tambahkan ke header response
		c.Writer.Header().Set("X-Request-ID", reqID)

		// Logging awal
		logrus.WithFields(logrus.Fields{
			"request_id": reqID,
			"ip":         ip,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
		}).Info("Incoming request")

		// Lanjut ke handler
		c.Next()

		// Logging selesai
		logrus.WithFields(logrus.Fields{
			"request_id": reqID,
			"ip":         ip,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   time.Since(start),
		}).Info("Request completed")
	}
}

func getClientIP(c *gin.Context) string {
	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.ClientIP()
	} else {
		ip = strings.Split(ip, ",")[0]
	}
	return strings.TrimSpace(ip)
}

// Get helper dari context Gin
func GetRequestID(c *gin.Context) string {
	if val, exists := c.Get(ContextKeyRequestID); exists {
		return val.(string)
	}
	return ""
}

func GetIPAddress(c *gin.Context) string {
	if val, exists := c.Get(ContextKeyIPAddress); exists {
		return val.(string)
	}
	return ""
}
