package route

import (
	"github.com/dwilanang/psp/internal/middleware"
	"github.com/dwilanang/psp/internal/registry"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, registry *registry.Registry) {
	h := registry.NewRoleHandler()

	rolesGroup := rg.Group("/roles")
	{
		rolesGroup.Use(middleware.RequireRole("SUPERADMIN"))
		rolesGroup.GET("/all", h.GetAll)
		rolesGroup.POST("/create", h.Create)
		rolesGroup.PUT("/update/:id", h.Update)
		rolesGroup.DELETE("/delete/:id", h.Delete)
	}
}
