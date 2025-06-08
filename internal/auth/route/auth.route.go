package route

import (
	"github.com/dwilanang/psp/internal/registry"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, registry *registry.Registry) {
	h := registry.NewAuthHandler()

	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
	}
}
