package main

import (
	"fmt"
	"log"

	"github.com/Thanhdat-debug/demo_login/internal/config"
	"github.com/Thanhdat-debug/demo_login/internal/handlers"
	"github.com/Thanhdat-debug/demo_login/internal/middleware"
	"github.com/Thanhdat-debug/demo_login/internal/models"
	"github.com/Thanhdat-debug/demo_login/internal/repository"
	"github.com/Thanhdat-debug/demo_login/internal/services"
	"github.com/Thanhdat-debug/demo_login/pkg/database"
	"github.com/Thanhdat-debug/demo_login/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Khởi tạo logger
	appLogger := logger.NewLogger()

	// Tải cấu hình từ file .env
	appConfig, err := config.LoadConfig()
	if err != nil {
		appLogger.Error("Failed to load config:", err)
		log.Fatal(err)
	}

	// Kết nối database
	db, err := database.NewDatabase(appConfig)
	if err != nil {
		appLogger.Error("Failed to connect to database:", err)
		log.Fatal(err)
	}

	// Khởi tạo gORM
	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Error("Failed to get database instance:", err)
		log.Fatal(err)
	}

	// Kiểm tra kết nối database
	if err := sqlDB.Ping(); err != nil {
		appLogger.Error("Failed to ping database:", err)
		log.Fatal(err)
	}
	appLogger.Info("Connected to database successfully")

	// Auto Migrate các model
	if err := db.AutoMigrate(&models.User{}); err != nil {
		appLogger.Error("Failed to auto migrate models:", err)
		log.Fatal(err)
	}
	appLogger.Info("Auto migration completed")

	// Khởi tạo repository
	userRepo := repository.NewUserRepository(db, appLogger)

	// Khởi tạo service
	authService := services.NewAuthService(userRepo, appConfig, appLogger)
	userService := services.NewUserService(userRepo, appLogger)

	// Khởi tạo middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, appLogger)

	// Khởi tạo handler
	authHandler := handlers.NewAuthHandler(authService, appLogger)
	userHandler := handlers.NewUserHandler(userService, appLogger)

	// Khởi tạo Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Định nghĩa các API route

	// Public routes (không cần xác thực)
	router.POST("/api/auth/register", authHandler.Register)
	router.POST("/api/auth/login", authHandler.Login)
	router.GET("/api/auth/validate", authHandler.ValidateToken)

	// Protected routes (cần xác thực JWT)
	protected := router.Group("/api")
	protected.Use(authMiddleware.JWTAuthMiddleware())
	{
		// User routes
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)
		protected.PUT("/users/change-password", userHandler.ChangePassword)
		protected.DELETE("/users/account", userHandler.DeleteAccount)

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(authMiddleware.AdminRequired())
		{
			admin.GET("/users", userHandler.GetUsersList)
		}
	}

	// Khởi động server
	serverAddr := fmt.Sprintf(":%s", appConfig.ServerPort)
	appLogger.Infof("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		appLogger.Error("Failed to start server:", err)
		log.Fatal(err)
	}
}
