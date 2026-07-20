package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"ai-forum/forum-service/internal/client"
	"ai-forum/forum-service/internal/handler"
	"ai-forum/forum-service/internal/middleware"
	"ai-forum/forum-service/internal/service"
	"ai-forum/forum-service/pkg/database"
	forumredis "ai-forum/forum-service/pkg/redis"

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

	// Set search path to service schema
	_, err = pool.Exec(context.Background(), fmt.Sprintf("SET search_path TO %s, public", dbSchema))
	if err != nil {
		log.Fatalf("Failed to set schema: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(dbHost, dbPort, dbName, dbUser, dbPassword); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	rdb, err := forumredis.Connect(redisHost, redisPort)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()
	log.Println("Connected to Redis")

	// Initialize admin service client
	adminURL := os.Getenv("ADMIN_SERVICE_URL")
	if adminURL == "" {
		adminURL = "http://localhost:8003"
	}
	adminClient := client.NewAdminClient(adminURL)

	userURL := os.Getenv("USER_SERVICE_URL")
	if userURL == "" {
		userURL = "http://localhost:8001"
	}
	userClient := client.NewUserClient(userURL)

	// Initialize service and handlers
	forumService := service.NewForumService(pool, adminClient, userClient)
	extrasService := service.NewExtrasService(pool)
	forumHandler := handler.NewForumHandler(forumService)
	commentHandler := handler.NewCommentHandler(forumService, extrasService)
	notificationHandler := handler.NewNotificationHandler(extrasService)
	reportHandler := handler.NewReportHandler(extrasService)
	reportAdminHandler := handler.NewReportAdminHandler(extrasService)
	communityHandler := handler.NewCommunityHandler(extrasService)
	interactionHandler := handler.NewInteractionHandler(forumService, extrasService)
	attachmentHandler := handler.NewAttachmentHandler(forumService)

	// Initialize forum admin service and handler
	forumAdminService := service.NewForumAdminService(pool)
	forumAdminHandler := handler.NewForumAdminHandler(forumAdminService)

	// Initialize stats service and handler
	statsService := service.NewStatsService(pool)
	statsHandler := handler.NewStatsHandler(statsService)

	// Initialize Gin router
	router := gin.Default()

	// Register routes
	router.GET("/health", handler.HealthCheck)

	// Public routes (no auth)
	v1 := router.Group("/api/v1")
	{
		// Board routes
		v1.GET("/boards", forumHandler.ListBoards)
		v1.GET("/boards/:id", forumHandler.GetBoard)

		// Post routes (read only; optional auth enriches liked/collected)
		postRead := v1.Group("")
		postRead.Use(middleware.OptionalAuthMiddleware())
		postRead.GET("/posts", forumHandler.ListPosts)
		postRead.GET("/posts/:id", forumHandler.GetPost)
		postRead.GET("/posts/:id/comments", commentHandler.ListComments)

		// Attachment download (read only)
		v1.GET("/attachments/:id", attachmentHandler.DownloadAttachment)
		v1.GET("/posts/:id/attachments", attachmentHandler.GetPostAttachments)
		v1.GET("/stats/community", communityHandler.Stats)
	}

	// Authenticated routes
	auth := v1.Group("")
	auth.Use(middleware.AuthMiddleware())
	// Sensitive words check (pre-flight for frontend)
	auth.POST("/moderation/check", func(c *gin.Context) {
		var req struct {
			Text string `json:"text" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
		clean, matched, err := adminClient.CheckSensitiveWords(req.Text)
		if err != nil {
			c.JSON(502, gin.H{"error": "moderation service unavailable"})
			return
		}
		c.JSON(200, gin.H{"clean": clean, "matched_words": matched})
	})
	{
		// Post routes (write)
		auth.POST("/posts", middleware.RequireLevel(0), middleware.RateLimitByUser(rdb, "create_post", 10, time.Hour), forumHandler.CreatePost)
		auth.PUT("/posts/:id", forumHandler.UpdatePost)
		auth.DELETE("/posts/:id", forumHandler.DeletePost)

		// Comment routes (write)
		auth.POST("/posts/:id/comments", middleware.RequireLevel(0), commentHandler.CreateComment)

		// Interaction routes
		auth.POST("/posts/:id/like", interactionHandler.LikePost)
		auth.POST("/posts/:id/dislike", interactionHandler.DislikePost)
		auth.POST("/comments/:id/like", interactionHandler.LikeComment)
		auth.POST("/comments/:id/dislike", interactionHandler.DislikeComment)
		auth.POST("/posts/:id/collect", interactionHandler.CollectPost)
		auth.GET("/me/collections", interactionHandler.ListMyCollections)

		// Attachment upload (requires level 2 for files)
		auth.POST("/attachments/upload", middleware.RequireLevel(2), attachmentHandler.UploadAttachment)

		auth.GET("/notifications", notificationHandler.List)
		auth.GET("/notifications/unread-count", notificationHandler.UnreadCount)
		auth.PUT("/notifications/:id/read", notificationHandler.MarkRead)
		auth.POST("/posts/:id/report", reportHandler.ReportPost)
	}

	// Internal routes (service-to-service only)
	internal := router.Group("/internal/v1")
	{
		internal.POST("/moderation/check", func(c *gin.Context) {
			var req struct {
				Text string `json:"text" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "invalid request"})
				return
			}
			clean, _, err := adminClient.CheckSensitiveWords(req.Text)
			if err != nil {
				c.JSON(502, gin.H{"error": "moderation service unavailable"})
				return
			}
			c.JSON(200, gin.H{"clean": clean})
		})

		// Admin post management
		internal.POST("/posts/:id/admin-delete", forumAdminHandler.AdminDeletePost)
		internal.POST("/posts/:id/admin-featured", forumAdminHandler.SetPostFeatured)
		internal.POST("/posts/:id/admin-pinned", forumAdminHandler.SetPostPinned)
		internal.POST("/posts/:id/admin-status", forumAdminHandler.AdminStatus)
		internal.GET("/posts/pending", forumAdminHandler.ListPendingPosts)
		internal.GET("/posts/all", forumAdminHandler.ListAllPosts)
		internal.POST("/posts/batch-delete", forumAdminHandler.BatchDeletePosts)

		// Admin board management
		internal.POST("/boards/admin", forumAdminHandler.CreateBoard)
		internal.PUT("/boards/admin/:id", forumAdminHandler.UpdateBoard)
		internal.DELETE("/boards/admin/:id", forumAdminHandler.DeleteBoard)
		internal.GET("/boards", forumAdminHandler.ListAllBoardsInternal)

		// Stats endpoints
		internal.GET("/stats/overview", statsHandler.GetStatsOverview)
		internal.GET("/stats/daily-posts", statsHandler.GetDailyPosts)
		internal.GET("/stats/daily-comments", statsHandler.GetDailyComments)
		internal.GET("/stats/board-activity", statsHandler.GetBoardActivity)

		internal.GET("/reports", reportAdminHandler.ListReports)
		internal.POST("/reports/:id/resolve", reportAdminHandler.ResolveReport)
	}

	// Start server
	port := ":8002"
	log.Printf("Starting forum-service on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
