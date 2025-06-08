package route

import (
	"github.com/dwilanang/psp/internal/middleware"
	"github.com/dwilanang/psp/internal/registry"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, registry *registry.Registry) {
	h := registry.NewAdminHandler()

	adminGroup := rg.Group("/admin")
	{
		adminGroup.Use(middleware.RequireRole("SUPERADMIN", "ADMIN"))
		adminGroup.POST("/attendance-periods", h.CreateAttendancePeriod)
		adminGroup.GET("/attendance-periods/all", h.FetchAttendancePeriod)
		adminGroup.POST("/payroll", h.CreatePayroll)
		adminGroup.GET("/payroll/all", h.FetchPayroll)
		adminGroup.PUT("/payroll/:id/process", h.RunPayroll)
		adminGroup.GET("/payroll/:id/summary", h.SummaryPayroll)
	}
}
