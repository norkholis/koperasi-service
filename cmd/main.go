package main

import (
	"koperasi-service/config"
	"koperasi-service/internal/handler"
	"koperasi-service/internal/middleware"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
	"koperasi-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPass + " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect db")
	}

	// Auto migrate
	db.AutoMigrate(&model.User{}, &model.Role{}, &model.Simpanan{})

	// Seed roles
	seedRoles(db)

	// Setup dependencies
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, cfg)

	// Additional services
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())

	// Public
	r.POST("/api/register", authHandler.Register)
	r.POST("/api/login", authHandler.Login)

	// Simpanan dependencies
	simpananRepo := repository.NewSimpananRepository(db)
	simpananSvc := service.NewSimpananService(simpananRepo)
	simpananHdl := handler.NewSimpananHandler(simpananSvc)

	// Protected
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg, userRepo))
	{
		protected.GET("/me", authHandler.Me)
		// Simpanan CRUD
		protected.POST("/simpanan", simpananHdl.Create)
		protected.GET("/simpanan", simpananHdl.List)
		protected.GET("/simpanan/:id", simpananHdl.Detail)
		protected.PUT("/simpanan/:id", simpananHdl.Update)
		protected.DELETE("/simpanan/:id", simpananHdl.Delete)

		// User CRUD
		protected.GET("/users", userHandler.List)
		protected.GET("/users/:id", userHandler.Detail)
		protected.POST("/users", userHandler.Create)
		protected.PUT("/users/:id", userHandler.Update)
		protected.DELETE("/users/:id", userHandler.Delete)
	}

	r.Run(":8080")
}

func seedRoles(db *gorm.DB) {
	roles := []model.Role{
		{Name: "super_admin"},
		{Name: "admin"},
		{Name: "member"},
	}
	for _, r := range roles {
		db.FirstOrCreate(&r, model.Role{Name: r.Name})
	}
}
