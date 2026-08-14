package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"ai-forum/admin-service/internal/handler"
	"ai-forum/admin-service/internal/middleware"
	"ai-forum/admin-service/internal/service"
	"ai-forum/admin-service/pkg/database"

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

	configService := service.NewConfigService(pool)
	userAdminService := service.NewUserAdminService(pool, userClient)
	roleService := service.NewRoleService(pool)

	inviteHandler := handler.NewInviteHandler(adminService, userClient)
	moderationHandler := handler.NewModerationHandler(adminService)
	configHandler := handler.NewConfigHandler(configService)
	userAdminHandler := handler.NewUserAdminHandler(userAdminService)
	inviteAdminHandler := handler.NewInviteAdminHandler(userClient)
	roleHandler := handler.NewRoleHandler(roleService)

	router := gin.Default()
	router.GET("/health", handler.HealthCheck)

	v1 := router.Group("/api/v1/admin")
	v1.Use(middleware.AuthMiddleware())
	{
		v1.GET("/config", configHandler.GetConfig)
		v1.GET("/config/:key", configHandler.GetConfigByKey)
		v1.PUT("/config/:key", configHandler.UpdateConfig)

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

		v1.POST("/sensitive-words", moderationHandler.AddSensitiveWord)
		v1.GET("/sensitive-words", moderationHandler.ListSensitiveWords)
		v1.DELETE("/sensitive-words/:id", moderationHandler.DeleteSensitiveWord)
	}

	internal := router.Group("/internal/v1")
	{
		internal.POST("/moderation/check", moderationHandler.CheckSensitiveWords)
		internal.GET("/users/:id/role", roleHandler.GetUserRole)
	}

	port := ":8003"
	log.Printf("Starting admin-service on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}