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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.DB.DSN()), &gorm.Config{
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
	qrRepo := repository.NewQRCodeRepository(db)
	queueRepo := repository.NewQueueRepository(db)

	// Services
	branchSvc := service.NewBranchService(branchRepo)
	qrSvc := service.NewQRCodeService(qrRepo, branchRepo, cfg.QR)
	queueSvc := service.NewQueueService(queueRepo, qrRepo, branchRepo)

	// Handlers
	branchHandler := handler.NewBranchHandler(branchSvc)
	qrHandler := handler.NewQRCodeHandler(qrSvc)
	queueHandler := handler.NewQueueHandler(queueSvc)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.LoadHTMLGlob("web/templates/*")
	r.Static("/storage/qr", cfg.QR.StoragePath)

	// Internal routes (protected by API key)
	internal := r.Group("/internal", middleware.InternalAuth(cfg.Internal.APIKey))
	{
		branches := internal.Group("/branches")
		{
			branches.POST("", branchHandler.Create)
			branches.GET("", branchHandler.List)
			branches.PUT("/:id", branchHandler.Update)
			branches.DELETE("/:id", branchHandler.Delete)
			branches.POST("/:id/qr", qrHandler.Generate)
			branches.POST("/:id/next", queueHandler.CallNext)
			branches.GET("/:id/queue", queueHandler.ListByBranch)
			branches.POST("/:id/reset", queueHandler.Reset)
		}
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
