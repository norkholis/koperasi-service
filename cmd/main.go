package main

import (
	"fmt"
	"log"
	"os"

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

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect db")
	}

	// Auto migrate
	db.AutoMigrate(&model.User{}, &model.Role{}, &model.Simpanan{}, &model.SimpananTransaction{}, &model.Pinjaman{}, &model.Angsuran{}, &model.SHUTahunan{}, &model.SHUAnggotaRecord{}, &model.BungaOption{}, &model.AuditTrail{}, &model.TransactionHistory{})

	// Seed roles
	seedRoles(db)

	// Setup dependencies
	userRepo := repository.NewUserRepository(db)
	simpananRepo := repository.NewSimpananRepository(db)

	// Auth service
	authService := service.NewAuthService(userRepo, simpananRepo)
	authHandler := handler.NewAuthHandler(authService, cfg)

	// User service
	userService := service.NewUserService(userRepo, simpananRepo)
	userHandler := handler.NewUserHandler(userService)

	// Basic repositories
	pinjamanRepo := repository.NewPinjamanRepository(db)
	angsuranRepo := repository.NewAngsuranRepository(db)

	// Audit Trail and Transaction History dependencies (need to be created first)
	auditRepo := repository.NewAuditTrailRepository(db)
	transactionRepo := repository.NewTransactionHistoryRepository(db)
	auditSvc := service.NewAuditTrailService(auditRepo, userRepo)
	transactionSvc := service.NewTransactionHistoryService(transactionRepo, userRepo)
	auditHdl := handler.NewAuditTrailHandler(auditSvc, transactionSvc)

	// Simpanan service (requires transaction service)
	simpananSvc := service.NewSimpananService(simpananRepo, transactionSvc)
	simpananHdl := handler.NewSimpananHandler(simpananSvc)

	// SHU dependencies
	shuRepo := repository.NewSHUTahunanRepository(db)
	shuSvc := service.NewSHUService(shuRepo)
	shuHdl := handler.NewSHUHandler(shuSvc)

	// SHU Anggota dependencies (requires transaction service)
	shuAnggotaRepo := repository.NewSHUAnggotaRepository(db)
	shuAnggotaSvc := service.NewSHUAnggotaService(shuAnggotaRepo, shuRepo, transactionSvc)
	shuAnggotaHdl := handler.NewSHUAnggotaHandler(shuAnggotaSvc)

	// Bunga Option dependencies
	bungaOptionRepo := repository.NewBungaOptionRepository(db)
	bungaOptionSvc := service.NewBungaOptionService(bungaOptionRepo, userRepo)
	bungaOptionHdl := handler.NewBungaOptionHandler(bungaOptionSvc)

	// Pinjaman service (requires transaction service)
	pinjamanSvc := service.NewPinjamanService(pinjamanRepo, userRepo, bungaOptionRepo, transactionSvc)
	pinjamanHdl := handler.NewPinjamanHandler(pinjamanSvc)

	// Angsuran service (requires transaction service)
	angsuranSvc := service.NewAngsuranService(angsuranRepo, pinjamanRepo, userRepo, transactionSvc)
	angsuranHdl := handler.NewAngsuranHandler(angsuranSvc)

	// Bank Account dependencies
	bankAccountRepo := repository.NewBankAccountRepository(db)
	bankAccountSvc := service.NewBankAccountService(bankAccountRepo)
	bankAccountHdl := handler.NewBankAccountHandler(bankAccountSvc)

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

	// Public routes
	r.POST("/api/register", authHandler.Register)
	r.POST("/api/login", authHandler.Login)
	r.POST("/api/forgot-password", authHandler.ForgotPassword)

	// Protected
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg, userRepo))
	{
		protected.GET("/me", authHandler.Me)
		protected.PUT("/me/profile", authHandler.UpdateProfile)
		protected.PUT("/me/password", authHandler.ChangePassword)
		protected.POST("/change-password", authHandler.ChangePassword)

		// User CRUD
		protected.GET("/users", userHandler.List)
		protected.POST("/users", userHandler.Create)
		protected.GET("/users/:id", userHandler.Detail)
		protected.PUT("/users/:id", userHandler.Update)
		protected.DELETE("/users/:id", userHandler.Delete)

		// Simpanan (Wallet) Management
		protected.GET("/simpanan/wallets", simpananHdl.GetWallets)
		protected.GET("/simpanan/wallets/all", simpananHdl.GetAllWallets)
		protected.POST("/simpanan/topup", simpananHdl.TopupWallet)
		protected.GET("/simpanan/:id", simpananHdl.GetWalletDetail)
		protected.GET("/simpanan/:id/transactions", simpananHdl.GetWalletTransactions)
		protected.PUT("/simpanan/:id/adjust", simpananHdl.AdjustWallet)
		protected.GET("/simpanan/transactions/pending", simpananHdl.GetPendingTransactions)
		protected.PUT("/simpanan/transactions/:id/verify", simpananHdl.VerifyTransaction)

		// Bank Account Management
		protected.POST("/bank-accounts", bankAccountHdl.Create)
		protected.GET("/bank-accounts", bankAccountHdl.List)
		protected.GET("/bank-accounts/active", bankAccountHdl.GetActive)
		protected.GET("/bank-accounts/:id", bankAccountHdl.GetByID)
		protected.PUT("/bank-accounts/:id", bankAccountHdl.Update)
		protected.DELETE("/bank-accounts/:id", bankAccountHdl.Delete)

		// Bunga Options Management
		protected.POST("/bunga-options", bungaOptionHdl.Create)
		protected.GET("/bunga-options", bungaOptionHdl.List)
		protected.GET("/bunga-options/:id", bungaOptionHdl.Detail)
		protected.PUT("/bunga-options/:id", bungaOptionHdl.Update)
		protected.DELETE("/bunga-options/:id", bungaOptionHdl.Delete)
		protected.PUT("/bunga-options/:id/status", bungaOptionHdl.SetActive)

		// Pinjaman CRUD with business logic
		protected.POST("/pinjaman", pinjamanHdl.Create)
		protected.GET("/pinjaman", pinjamanHdl.List)
		protected.GET("/pinjaman/:id", pinjamanHdl.Detail)
		protected.PUT("/pinjaman/:id", pinjamanHdl.Update)
		protected.DELETE("/pinjaman/:id", pinjamanHdl.Delete)
		protected.PUT("/pinjaman/:id/approve", pinjamanHdl.Approve)
		protected.PUT("/pinjaman/:id/reject", pinjamanHdl.Reject)

		// Angsuran CRUD with verification support
		protected.POST("/angsuran", angsuranHdl.Create)
		protected.GET("/angsuran", angsuranHdl.List)
		protected.GET("/angsuran/:id", angsuranHdl.Detail)
		protected.PUT("/angsuran/:id", angsuranHdl.Update)
		protected.DELETE("/angsuran/:id", angsuranHdl.Delete)
		protected.PUT("/angsuran/:id/verify", angsuranHdl.Verify)
		protected.GET("/angsuran/pending", angsuranHdl.GetPendingPayments)

		// SHU Management
		protected.POST("/shu/generate", shuHdl.GenerateReport)
		protected.POST("/shu/generate-auto", shuHdl.GenerateReportWithExpenses)
		protected.POST("/shu/save-auto", shuHdl.SaveSHUWithExpenses)
		protected.POST("/shu", shuHdl.SaveSHU)
		protected.GET("/shu", shuHdl.List)
		protected.GET("/shu/:id", shuHdl.Detail)
		protected.PUT("/shu/:id", shuHdl.Update)
		protected.DELETE("/shu/:id", shuHdl.Delete)
		protected.GET("/shu/year/:tahun", shuHdl.GetByTahun)
		protected.POST("/shu/user/:user_id/generate", shuHdl.GenerateUserSHU)

		// SHU Anggota Management
		protected.POST("/shu-anggota/user/:user_id/save", shuAnggotaHdl.SaveUserSHU)
		protected.GET("/shu-anggota/user/:user_id/:tahun", shuAnggotaHdl.GetUserSHU)
		protected.GET("/shu-anggota/user/:user_id/history", shuAnggotaHdl.GetUserSHUHistory)
		protected.GET("/shu-anggota", shuAnggotaHdl.List)
		protected.GET("/shu-anggota/shu/:shu_id", shuAnggotaHdl.GetBySHUID)
		protected.DELETE("/shu-anggota/:id", shuAnggotaHdl.Delete)

		// Transaction History & Financial Reporting
		protected.GET("/transactions", auditHdl.GetTransactionHistory)
		protected.GET("/transactions/:id", auditHdl.GetTransactionDetail)
		protected.GET("/transactions/user/:user_id", auditHdl.GetUserTransactions)
		protected.GET("/transactions/summary", auditHdl.GetFinancialSummary)
		protected.POST("/reports/financial", auditHdl.GenerateFinancialReport)
		protected.GET("/reports/financial", auditHdl.GenerateFinancialReportGET)

		// Audit Trail & System Monitoring
		protected.GET("/audit-trails", auditHdl.GetAuditTrails)
		protected.GET("/audit-trails/:id", auditHdl.GetAuditTrailDetail)
		protected.GET("/audit-trails/user/:user_id", auditHdl.GetUserActivity)
		protected.GET("/audit-trails/system", auditHdl.GetSystemActivity)
		protected.GET("/audit-trails/summary", auditHdl.GetAuditSummary)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}

func seedRoles(db *gorm.DB) {
	var roleCount int64
	db.Model(&model.Role{}).Count(&roleCount)

	if roleCount == 0 {
		roles := []model.Role{
			{Name: "super_admin"},
			{Name: "admin"},
			{Name: "member"},
		}

		for _, role := range roles {
			db.Create(&role)
		}

		log.Println("Roles seeded successfully")
	}
}
