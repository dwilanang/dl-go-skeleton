package route

import (
	"github.com/dwilanang/psp/internal/middleware"
	"github.com/dwilanang/psp/internal/registry"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, registry *registry.Registry) {
	h := registry.NewUserHandler()

	usersGroup := rg.Group("/users")
	{
		usersGroup.Use(middleware.RequireRole("SUPERADMIN"))
		usersGroup.POST("/register", h.Register)
	}
}
