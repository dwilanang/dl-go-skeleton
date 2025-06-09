package main

import (
	"fmt"

	"github.com/dwilanang/psp/config"
	_ "github.com/dwilanang/psp/docs"
	"github.com/dwilanang/psp/infrastructure/db/postgres"

	authroute "github.com/dwilanang/psp/internal/auth/route"
	"github.com/dwilanang/psp/internal/middleware"
	"github.com/dwilanang/psp/internal/registry"
	roleroute "github.com/dwilanang/psp/internal/role/route"
	userroute "github.com/dwilanang/psp/internal/user/route"
	"github.com/dwilanang/psp/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GO SKELETON API
// @version 1.0
// @description go skeleton project API
// @host localhost:8000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cfg := config.LoadConfig()
	// This is the entry point for the PSP application.
	// You can initialize your application here, set up routes, etc.

	// Initialize a database postgres connection
	dbPostgres := postgres.Connect(cfg)
	if dbPostgres == nil {
		fmt.Println("Failed to initialize the database connection")
		return
	}

	r := gin.Default()
	r.Use(logger.RequestLogger())

	registry := registry.NewRegistry(cfg, dbPostgres)

	api := r.Group("/api/v1")

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Auth
	authroute.RegisterRoutes(api, registry)

	api.Use(middleware.JWTAuthMiddleware(cfg.JWTSecret))

	// Register all feature routes here, keep it simple
	roleroute.RegisterRoutes(api, registry)
	userroute.RegisterRoutes(api, registry)

	fmt.Println("Server is running on port ", cfg.AppPort)
	r.Run(fmt.Sprintf(":%s", cfg.AppPort))
}
