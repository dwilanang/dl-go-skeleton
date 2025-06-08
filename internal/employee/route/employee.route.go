package route

import (
	"github.com/dwilanang/psp/internal/middleware"
	"github.com/dwilanang/psp/internal/registry"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, registry *registry.Registry) {
	h := registry.NewEmployeeHandler()

	employeeGroup := rg.Group("/employee")
	{
		employeeGroup.Use(middleware.RequireRole("SUPERADMIN", "EMPLOYEE"))
		employeeGroup.POST("/attendance", h.CreateAttendance)
		employeeGroup.GET("/attendance/all", h.FetchAttendance)
		employeeGroup.POST("/overtime", h.CreateOvertime)
		employeeGroup.GET("/overtime/all", h.FetchOvertime)
		employeeGroup.POST("/reimbursement", h.CreateReimbursement)
		employeeGroup.GET("/reimbursement/all", h.FetchReimbursement)
		employeeGroup.GET("/payslip/:id", h.GetEmployeePayslip)
	}
}
