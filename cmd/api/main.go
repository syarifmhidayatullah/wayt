package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/project/wayt/config"
	"github.com/project/wayt/internal/handler"
	"github.com/project/wayt/internal/repository"
	"github.com/project/wayt/internal/service"
	"github.com/project/wayt/pkg/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := os.MkdirAll(cfg.QR.StoragePath, 0755); err != nil {
		log.Fatalf("failed to create QR storage dir: %v", err)
	}

	// Repositories
	branchRepo := repository.NewBranchRepository(db)
	counterRepo := repository.NewCounterRepository(db)
	qrRepo := repository.NewQRCodeRepository(db)
	queueRepo := repository.NewQueueRepository(db)
	adminRepo := repository.NewAdminUserRepository(db)

	// Services
	branchSvc := service.NewBranchService(branchRepo)
	counterSvc := service.NewCounterService(counterRepo, branchRepo)
	qrSvc := service.NewQRCodeService(qrRepo, counterRepo, cfg.QR)
	queueSvc := service.NewQueueService(queueRepo, qrRepo, counterRepo)
	authSvc := service.NewAuthService(adminRepo, cfg.Auth.JWTSecret)

	// Seed default admin jika belum ada
	if cfg.Auth.AdminPassword != "" {
		if err := authSvc.SeedAdmin(cfg.Auth.AdminUsername, cfg.Auth.AdminPassword); err != nil {
			log.Printf("seed admin skipped: %v", err)
		}
	}

	// Handlers
	branchHandler := handler.NewBranchHandler(branchSvc)
	counterHandler := handler.NewCounterHandler(counterSvc, qrSvc, queueSvc)
	queueHandler := handler.NewQueueHandler(queueSvc)
	authHandler := handler.NewAuthHandler(authSvc)
	adminUserHandler := handler.NewAdminUserHandler(authSvc)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.LoadHTMLGlob("web/templates/*")
	r.Static("/storage/qr", cfg.QR.StoragePath)

	// Auth routes (public)
	r.POST("/auth/login", authHandler.Login)

	// Internal routes (protected by JWT)
	internal := r.Group("/internal", middleware.JWTAuth(cfg.Auth.JWTSecret))
	{
		// User management — superadmin only
		users := internal.Group("/users", middleware.SuperAdminOnly())
		{
			users.GET("", adminUserHandler.List)
			users.POST("", adminUserHandler.Create)
			users.PUT("/:id", adminUserHandler.Update)
			users.DELETE("/:id", adminUserHandler.Delete)
		}

		// Branch management — superadmin only for create/update/delete
		// List is available to all (filtered by role in handler)
		internal.GET("/branches", branchHandler.List)
		superBranches := internal.Group("/branches", middleware.SuperAdminOnly())
		{
			superBranches.POST("", branchHandler.Create)
			superBranches.PUT("/:id", branchHandler.Update)
			superBranches.DELETE("/:id", branchHandler.Delete)
		}

		// Counter management — admin restricted to their branch (checked in handler)
		internal.POST("/branches/:branch_id/counters", counterHandler.Create)
		internal.GET("/branches/:branch_id/counters", counterHandler.ListByBranch)
		internal.PUT("/counters/:id", counterHandler.Update)
		internal.DELETE("/counters/:id", counterHandler.Delete)
		internal.POST("/counters/:id/qr", counterHandler.GenerateQR)
		internal.POST("/counters/:id/next", counterHandler.CallNext)
		internal.GET("/counters/:id/queue", counterHandler.ListQueue)
		internal.POST("/counters/:id/reset", counterHandler.Reset)
	}

	// Public routes
	api := r.Group("/api")
	{
		api.POST("/queue/register", queueHandler.Register)
		api.GET("/queue/:token/status", queueHandler.Status)
		api.GET("/queue/id/:id/status", queueHandler.StatusByID)
	}

	// QR scan route — register lalu redirect ke /queue/:id
	r.GET("/q/:token", queueHandler.ScanRegister)

	// Halaman status antrian per orang (unique per queue entry)
	r.GET("/queue/:id", queueHandler.QueuePage)

	// Admin page
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "admin.html", gin.H{})
	})

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
