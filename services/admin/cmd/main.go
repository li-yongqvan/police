package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"ai-forum/admin-service/internal/client"
	"ai-forum/admin-service/internal/handler"
	"ai-forum/admin-service/internal/middleware"
	"ai-forum/admin-service/internal/service"
	"ai-forum/admin-service/pkg/database"
	"ai-forum/admin-service/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

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

	_, err = pool.Exec(context.Background(), fmt.Sprintf("SET search_path TO %s, public", dbSchema))
	if err != nil {
		log.Fatalf("Failed to set schema: %v", err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	rdb, err := redis.Connect(redisHost, redisPort)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()
	log.Println("Connected to Redis")

	if err := database.RunMigrations(dbHost, dbPort, dbName, dbUser, dbPassword); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	adminService := service.NewAdminService(pool)
	if err := adminService.LoadSensitiveWords(context.Background()); err != nil {
		log.Printf("Warning: failed to load sensitive words cache: %v", err)
	}

	userURL := os.Getenv("USER_SERVICE_URL")
	if userURL == "" {
		userURL = "http://localhost:8001"
	}
	userClient := service.NewUserClient(userURL)

	forumURL := os.Getenv("FORUM_SERVICE_URL")
	if forumURL == "" {
		forumURL = "http://localhost:8002"
	}
	forumClient := client.NewForumClient(forumURL)

	auditService := service.NewAuditService(pool, forumClient)
	configService := service.NewConfigService(pool)
	userAdminService := service.NewUserAdminService(pool, userClient)
	roleService := service.NewRoleService(pool)
	roleService.SetRedis(rdb)
	postAdminService := service.NewPostAdminService(forumClient, pool)
	statsService := service.NewStatsService(userClient, forumClient, pool)

	inviteHandler := handler.NewInviteHandler(adminService, userClient)
	moderationHandler := handler.NewModerationHandler(adminService)
	auditHandler := handler.NewAuditHandler(auditService)
	configHandler := handler.NewConfigHandler(configService)
	userAdminHandler := handler.NewUserAdminHandler(userAdminService)
	inviteAdminHandler := handler.NewInviteAdminHandler(userClient)
	boardAdminHandler := handler.NewBoardAdminHandler(forumClient)
	roleHandler := handler.NewRoleHandler(roleService)
	postAdminHandler := handler.NewPostAdminHandler(postAdminService, forumClient)
	statsHandler := handler.NewStatsHandler(statsService)

	router := gin.Default()
	router.GET("/health", handler.HealthCheck)

	v1 := router.Group("/api/v1/admin")
	v1.Use(middleware.AuthMiddleware())
	{
		v1.GET("/config", configHandler.GetConfig)
		v1.GET("/config/:key", configHandler.GetConfigByKey)
		v1.PUT("/config/:key", configHandler.UpdateConfig)

		v1.GET("/audit/pending", auditHandler.ListPendingAudit)
		v1.POST("/audit/:id/approve", auditHandler.ApprovePost)
		v1.POST("/audit/:id/reject", auditHandler.RejectPost)
		v1.POST("/audit/batch-delete", auditHandler.BatchDeletePosts)

		v1.POST("/posts/:id/delete", postAdminHandler.DeletePost)
		v1.GET("/posts", postAdminHandler.ListPosts)
		v1.POST("/posts/:id/featured", postAdminHandler.SetPostFeatured)
		v1.POST("/posts/:id/pinned", postAdminHandler.SetPostPinned)

		v1.POST("/users/:id/ban", userAdminHandler.BanUser)
		v1.POST("/users/:id/unban", userAdminHandler.UnbanUser)
		v1.GET("/users", userAdminHandler.ListUsers)
		v1.PUT("/users/:id/level", userAdminHandler.UpdateUserLevel)
		v1.GET("/users/:id/logs", userAdminHandler.GetUserLogs)

		v1.GET("/roles", roleHandler.ListRoles)
		v1.POST("/users/:id/roles", roleHandler.AssignRole)
		v1.DELETE("/users/:id/roles/:role_id", roleHandler.RemoveRole)
		v1.GET("/users/:id/roles", roleHandler.GetUserRoles)

		v1.POST("/invite-codes", inviteHandler.GenerateInviteCode)
		v1.POST("/invite-codes/batch", inviteHandler.GenerateInviteCodesBatch)
		v1.GET("/invite-codes", inviteAdminHandler.ListInviteCodes)
		v1.GET("/invite-codes/:code/status", inviteAdminHandler.GetInviteCodeStatus)
		v1.PUT("/invite-codes/:code/void", inviteAdminHandler.VoidInviteCode)

		v1.POST("/boards", boardAdminHandler.CreateBoard)
		v1.PUT("/boards/:id", boardAdminHandler.UpdateBoard)
		v1.DELETE("/boards/:id", boardAdminHandler.DeleteBoard)
		v1.GET("/boards", boardAdminHandler.ListBoards)

		v1.POST("/sensitive-words", moderationHandler.AddSensitiveWord)
		v1.GET("/sensitive-words", moderationHandler.ListSensitiveWords)
		v1.DELETE("/sensitive-words/:id", moderationHandler.DeleteSensitiveWord)

		v1.GET("/stats/overview", statsHandler.GetOverview)
		v1.GET("/stats/daily", statsHandler.GetDailyStats)
	}

	internal := router.Group("/internal/v1")
	{
		internal.POST("/moderation/check", moderationHandler.CheckSensitiveWords)
	}

	port := ":8003"
	log.Printf("Starting admin-service on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
