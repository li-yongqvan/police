package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"ai-forum/user-service/internal/handler"
	"ai-forum/user-service/internal/middleware"
	"ai-forum/user-service/internal/service"
	"ai-forum/user-service/pkg/database"
	"ai-forum/user-service/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	// Connect to PostgreSQL
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSchema := os.Getenv("DB_SCHEMA")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	pool, err := database.Connect(dbHost, dbPort, dbName, dbUser, dbPassword)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to PostgreSQL")

	// Set pool for middleware access
	database.SetPool(pool)

	// Set search path to service schema
	_, err = pool.Exec(context.Background(), fmt.Sprintf("SET search_path TO %s, public", dbSchema))
	if err != nil {
		log.Fatalf("Failed to set schema: %v", err)
	}

	// Connect to Redis
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	rdb, err := redis.Connect(redisHost, redisPort)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()
	log.Println("Connected to Redis")

	// Run migrations
	if err := database.RunMigrations(dbHost, dbPort, dbName, dbUser, dbPassword); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	// Initialize service and handler
	userService := service.NewUserService(pool, rdb)
	userHandler := handler.NewUserHandler(userService)

	// Initialize user admin service and handler
	userAdminService := service.NewUserAdminService(pool, rdb)
	userAdminHandler := handler.NewUserAdminHandler(userAdminService)

	// Initialize stats service and handler
	userStatsService := service.NewUserStatsService(pool)
	userStatsHandler := handler.NewUserStatsHandler(userStatsService)

	// Initialize Gin router
	router := gin.Default()

	// Register routes
	router.GET("/health", handler.HealthCheck)

	// Public routes (no auth)
	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", middleware.RateLimitByIP(rdb, "login", 5, time.Minute), userHandler.Login)
		v1.POST("/demo-login", middleware.DemoLoginGuard(), userHandler.DemoLogin)
		v1.POST("/auth/refresh", userHandler.RefreshToken)
	}

	// Authenticated routes
	auth := v1.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/users/me", userHandler.Me)
		auth.GET("/users/:id", userHandler.GetProfile)
		auth.PUT("/users/:id", userHandler.UpdateProfile)
		auth.POST("/users/:id/avatar", userHandler.UploadAvatar)
	}

	// Internal routes (service-to-service, no external auth per D-09)
	internal := router.Group("/internal/v1")
	{
		internal.GET("/users/:id/status", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
				return
			}
			uid, uname, level, status, err := userService.GetUserStatus(c.Request.Context(), uint(id))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"user_id":  uid,
				"username": uname,
				"level":    level,
				"status":   status,
			})
		})

		internal.POST("/users/batch-status", func(c *gin.Context) {
			var req struct {
				UserIDs []uint `json:"user_ids"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			result := make(map[string]interface{})
			for _, uid := range req.UserIDs {
				_, uname, level, status, _ := userService.GetUserStatus(c.Request.Context(), uid)
				result[strconv.FormatUint(uint64(uid), 10)] = gin.H{
					"username": uname,
					"level":    level,
					"status":   status,
				}
			}
			c.JSON(http.StatusOK, result)
		})

		internal.POST("/invite-codes", func(c *gin.Context) {
			createdBy := int64(0) // Admin system
			code, err := userService.CreateInviteCode(c.Request.Context(), createdBy)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": code})
		})

		internal.POST("/invite-codes/batch", func(c *gin.Context) {
			var req struct {
				Count     int `json:"count"`
				CreatedBy int64
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			if req.Count <= 0 || req.Count > 1000 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "count must be between 1 and 1000"})
				return
			}
			codes, err := userService.CreateInviteCodesBatch(c.Request.Context(), req.Count, req.CreatedBy)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"codes": codes, "count": len(codes)})
		})

		// Admin user management
		internal.POST("/users/:id/ban", userAdminHandler.BanUser)
		internal.POST("/users/:id/unban", userAdminHandler.UnbanUser)
		internal.GET("/users", userAdminHandler.ListUsers)
		internal.PUT("/users/:id/level", userAdminHandler.UpdateUserLevel)
		internal.GET("/users/:id/logs", userAdminHandler.GetUserLogs)

		// Admin invite code management
		internal.GET("/invite-codes", userAdminHandler.ListInviteCodes)
		internal.GET("/invite-codes/:code/status", userAdminHandler.GetInviteCodeStatus)
		internal.PUT("/invite-codes/:code/void", userAdminHandler.VoidInviteCode)

			// Stats endpoints
			internal.GET("/stats/overview", userStatsHandler.GetStatsOverview)
			internal.GET("/stats/daily-users", userStatsHandler.GetDailyUsers)
			internal.GET("/stats/level-distribution", userStatsHandler.GetLevelDistribution)
	}

	// Start server
	port := ":8001"
	log.Printf("Starting user-service on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
