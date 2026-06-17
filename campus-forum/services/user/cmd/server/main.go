package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Bio        string `json:"bio"`
	Status     string `json:"status"`
}

type demoLoginRequest struct {
	Role string `json:"role"`
}

type banRequest struct {
	Status string `json:"status"`
}

type statsOverview struct {
	UserCount int `json:"userCount"`
}

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var (
	sessionMu sync.Mutex
	sessions  = map[string]string{}
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors())

	router.POST("/api/v1/demo-login", handleDemoLogin)
	router.GET("/api/v1/users/me", handleCurrentUser)
	router.GET("/api/v1/admin/users", handleListUsers)
	router.GET("/internal/v1/users/:id/status", handleUserStatus)
	router.POST("/api/v1/admin/users/:id/ban-sync", handleBanSync)
	router.GET("/internal/v1/stats/overview", handleStatsOverview)

	port := envOr("USER_SERVICE_PORT", "8001")
	router.Run(":" + port)
}

func handleDemoLogin(c *gin.Context) {
	var req demoLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response{Success: false, Message: "登录参数不完整"})
		return
	}

	users, err := readUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Success: false, Message: err.Error()})
		return
	}

	for _, item := range users {
		if item.Role == req.Role {
			token := fmt.Sprintf("demo-%s-%d", item.ID, time.Now().UnixNano())
			sessionMu.Lock()
			sessions[token] = item.ID
			sessionMu.Unlock()
			c.JSON(http.StatusOK, response{Success: true, Message: "登录成功", Data: gin.H{
				"token": token,
				"user":  item,
			}})
			return
		}
	}

	c.JSON(http.StatusNotFound, response{Success: false, Message: "未找到对应角色账号"})
}

func handleCurrentUser(c *gin.Context) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, response{Success: false, Message: "缺少登录态"})
		return
	}

	sessionMu.Lock()
	userID := sessions[token]
	sessionMu.Unlock()
	if userID == "" {
		c.JSON(http.StatusUnauthorized, response{Success: false, Message: "登录态已失效"})
		return
	}

	current, err := findUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response{Success: true, Message: "获取成功", Data: current})
}

func handleListUsers(c *gin.Context) {
	users, err := readUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Success: true, Message: "获取成功", Data: users})
}

func handleUserStatus(c *gin.Context) {
	item, err := findUser(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Success: true, Message: "获取成功", Data: item})
}

func handleBanSync(c *gin.Context) {
	var req banRequest
	if err := c.ShouldBindJSON(&req); err != nil || (req.Status != "active" && req.Status != "banned") {
		c.JSON(http.StatusBadRequest, response{Success: false, Message: "状态参数错误"})
		return
	}

	users, err := readUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Success: false, Message: err.Error()})
		return
	}

	targetID := c.Param("id")
	updated := false
	for i := range users {
		if users[i].ID == targetID {
			users[i].Status = req.Status
			updated = true
		}
	}
	if !updated {
		c.JSON(http.StatusNotFound, response{Success: false, Message: "用户不存在"})
		return
	}

	if err := writeUsers(users); err != nil {
		c.JSON(http.StatusInternalServerError, response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Success: true, Message: "用户状态已更新"})
}

func handleStatsOverview(c *gin.Context) {
	users, err := readUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Success: true, Message: "获取成功", Data: statsOverview{UserCount: len(users)}})
}

func findUser(id string) (user, error) {
	users, err := readUsers()
	if err != nil {
		return user{}, err
	}
	for _, item := range users {
		if item.ID == id {
			return item, nil
		}
	}
	return user{}, errors.New("用户不存在")
}

func readUsers() ([]user, error) {
	path := filepath.Join(mockDataDir(), "users.json")
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var users []user
	if err := json.Unmarshal(bytes, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func writeUsers(users []user) error {
	path := filepath.Join(mockDataDir(), "users.json")
	bytes, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

func mockDataDir() string {
	if custom := os.Getenv("MOCK_DATA_DIR"); custom != "" {
		return custom
	}
	return filepath.Clean(filepath.Join("..", "..", "shared", "mock-data"))
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
