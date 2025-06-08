package route

import (
	"github.com/dwilanang/psp/internal/middleware"
	"github.com/dwilanang/psp/internal/registry"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, registry *registry.Registry) {
	h := registry.NewAuditLogsHandler()

	employeeGroup := rg.Group("/auditlogs")
	{
		employeeGroup.Use(middleware.RequireRole("SUPERADMIN", "ADMIN"))
		employeeGroup.GET("/all", h.GetAuditLogs)
	}
}
