package main

import (
	"koperasi-service/config"
	"koperasi-service/internal/handler"
	"koperasi-service/internal/middleware"
	"koperasi-service/internal/model"
	"koperasi-service/internal/repository"
	"koperasi-service/internal/service"

	"github.com/gin-contrib/cors"
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
	db.AutoMigrate(&model.User{}, &model.Role{}, &model.Simpanan{}, &model.SimpananTransaction{}, &model.Pinjaman{}, &model.Angsuran{}, &model.SHUTahunan{})

	// Seed roles
	seedRoles(db)

	// Setup dependencies
	userRepo := repository.NewUserRepository(db)
	simpananRepo := repository.NewSimpananRepository(db)

	authService := service.NewAuthService(userRepo, simpananRepo)
	authHandler := handler.NewAuthHandler(authService, cfg)

	// Additional services
	userService := service.NewUserService(userRepo, simpananRepo)
	userHandler := handler.NewUserHandler(userService)

	// Pinjaman dependencies
	pinjamanRepo := repository.NewPinjamanRepository(db)
	pinjamanSvc := service.NewPinjamanService(pinjamanRepo)
	pinjamanHdl := handler.NewPinjamanHandler(pinjamanSvc)

	// Angsuran dependencies
	angsuranRepo := repository.NewAngsuranRepository(db)
	angsuranSvc := service.NewAngsuranService(angsuranRepo, pinjamanRepo)
	angsuranHdl := handler.NewAngsuranHandler(angsuranSvc)

	// SHU dependencies
	shuRepo := repository.NewSHUTahunanRepository(db)
	shuSvc := service.NewSHUService(shuRepo)
	shuHdl := handler.NewSHUHandler(shuSvc)

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins, change in production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(middleware.LoggerMiddleware())

	// Public
	r.POST("/api/register", authHandler.Register)
	r.POST("/api/login", authHandler.Login)

	// Simpanan dependencies
	simpananSvc := service.NewSimpananService(simpananRepo)
	simpananHdl := handler.NewSimpananHandler(simpananSvc)

	// Protected
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg, userRepo))
	{
		protected.GET("/me", authHandler.Me)
		// Simpanan (Wallet) Management
		protected.GET("/simpanan/wallets", simpananHdl.GetWallets)                          // Get user wallets (with optional user_id param for admin)
		protected.GET("/simpanan/wallets/all", simpananHdl.GetAllWallets)                   // Get all wallets (admin only)
		protected.POST("/simpanan/topup", simpananHdl.TopupWallet)                          // Top-up wallet (creates pending transaction)
		protected.GET("/simpanan/:id", simpananHdl.GetWalletDetail)                         // Get wallet detail
		protected.GET("/simpanan/:id/transactions", simpananHdl.GetWalletTransactions)      // Get wallet transaction history
		protected.PUT("/simpanan/:id/adjust", simpananHdl.AdjustWallet)                     // Admin adjust wallet balance
		protected.GET("/simpanan/transactions/pending", simpananHdl.GetPendingTransactions) // Get pending transactions (admin)
		protected.PUT("/simpanan/transactions/:id/verify", simpananHdl.VerifyTransaction)   // Verify transaction (admin)

		// User CRUD
		protected.GET("/users", userHandler.List)
		protected.GET("/users/:id", userHandler.Detail)
		protected.POST("/users", userHandler.Create)
		protected.PUT("/users/:id", userHandler.Update)
		protected.DELETE("/users/:id", userHandler.Delete)

		// Pinjaman CRUD
		protected.GET("/pinjaman", pinjamanHdl.List)
		protected.GET("/pinjaman/:id", pinjamanHdl.Detail)
		protected.POST("/pinjaman", pinjamanHdl.Create)
		protected.PUT("/pinjaman/:id", pinjamanHdl.Update)
		protected.DELETE("/pinjaman/:id", pinjamanHdl.Delete)

		// Angsuran CRUD
		protected.GET("/angsuran", angsuranHdl.List)
		protected.GET("/angsuran/:id", angsuranHdl.Detail)
		protected.POST("/angsuran", angsuranHdl.Create)
		protected.PUT("/angsuran/:id", angsuranHdl.Update)
		protected.DELETE("/angsuran/:id", angsuranHdl.Delete)
		protected.PUT("/angsuran/:id/verify", angsuranHdl.Verify)
		protected.GET("/angsuran/pending", angsuranHdl.GetPendingPayments)

		// SHU (Sisa Hasil Usaha) - Admin/Super Admin only
		protected.POST("/shu/generate", shuHdl.GenerateReport)
		protected.POST("/shu", shuHdl.SaveSHU)
		protected.GET("/shu", shuHdl.List)
		protected.GET("/shu/:id", shuHdl.Detail)
		protected.PUT("/shu/:id", shuHdl.Update)
		protected.DELETE("/shu/:id", shuHdl.Delete)
		protected.GET("/shu/year/:tahun", shuHdl.GetByTahun)
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
