package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type postCheck struct {
	ID       string `json:"id"`
	BoardID  string `json:"boardId"`
	AuthorID string `json:"authorId"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type auditRecord struct {
	ID         string `json:"id"`
	PostID     string `json:"postId"`
	Reason     string `json:"reason"`
	Status     string `json:"status"`
	ReviewerID string `json:"reviewerId"`
	CreatedAt  string `json:"createdAt"`
}

type systemConfig struct {
	PostingEnabled bool            `json:"postingEnabled"`
	BoardSwitches  map[string]bool `json:"boardSwitches"`
	ModerationMode string          `json:"moderationMode"`
}

type actionRequest struct {
	ReviewerID string `json:"reviewerId"`
}

type configRequest struct {
	PostingEnabled bool            `json:"postingEnabled"`
	BoardSwitches  map[string]bool `json:"boardSwitches"`
	ModerationMode string          `json:"moderationMode"`
}

type banRequest struct {
	Status string `json:"status"`
}

type apiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors())

	router.GET("/api/v1/admin/config", handleGetConfig)
	router.PUT("/api/v1/admin/config", handlePutConfig)
	router.GET("/api/v1/admin/audit/pending", handlePendingAudit)
	router.POST("/api/v1/admin/audit/:id/approve", handleApprove)
	router.POST("/api/v1/admin/audit/:id/reject", handleReject)
	router.POST("/api/v1/admin/posts/:id/delete", handleDeletePost)
	router.POST("/api/v1/admin/users/:id/ban", handleBanUser)
	router.GET("/api/v1/admin/stats/overview", handleStatsOverview)
	router.POST("/internal/v1/moderation/check-post", handleModerationCheck)

	port := envOr("ADMIN_SERVICE_PORT", "8003")
	router.Run(":" + port)
}

func handleGetConfig(c *gin.Context) {
	cfg, err := readConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "获取成功", Data: cfg})
}

func handlePutConfig(c *gin.Context) {
	var req configRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Message: "配置参数错误"})
		return
	}
	cfg := systemConfig(req)
	if cfg.ModerationMode == "" {
		cfg.ModerationMode = "auto"
	}
	if err := writeJSON("system-config.json", cfg); err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "配置已更新", Data: cfg})
}

func handleModerationCheck(c *gin.Context) {
	var req postCheck
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Message: "审核输入错误"})
		return
	}
	cfg, err := readConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	sensitiveWords, _ := readSensitiveWords()
	reason := ""
	status := "published"
	lowerText := strings.ToLower(req.Title + " " + req.Content)
	for _, word := range sensitiveWords {
		if strings.Contains(lowerText, strings.ToLower(word)) {
			reason = "命中敏感词：" + word
			status = "pending"
			break
		}
	}
	if reason == "" && cfg.ModerationMode == "manual" {
		reason = "当前为人工审核模式"
		status = "pending"
	}
	if status == "pending" {
		records, _ := readAuditRecords()
		records = append([]auditRecord{{
			ID:         "audit-" + time.Now().Format("20060102150405"),
			PostID:     req.ID,
			Reason:     reason,
			Status:     "pending",
			ReviewerID: "",
			CreatedAt:  time.Now().Format(time.RFC3339),
		}}, records...)
		if err := writeJSON("audit-records.json", records); err != nil {
			c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "审核完成", Data: gin.H{
		"status": status,
		"reason": reason,
	}})
}

func handlePendingAudit(c *gin.Context) {
	records, err := readAuditRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	posts, _ := readPosts()
	postMap := map[string]map[string]string{}
	for _, item := range posts {
		postMap[item["id"]] = item
	}
	items := []gin.H{}
	for _, record := range records {
		if record.Status == "pending" {
			meta := postMap[record.PostID]
			items = append(items, gin.H{
				"id":        record.ID,
				"postId":    record.PostID,
				"reason":    record.Reason,
				"createdAt": record.CreatedAt,
				"title":     meta["title"],
				"authorId":  meta["authorId"],
				"boardId":   meta["boardId"],
			})
		}
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "获取成功", Data: items})
}

func handleApprove(c *gin.Context) {
	reviewAction(c, "approve")
}

func handleReject(c *gin.Context) {
	reviewAction(c, "reject")
}

func handleDeletePost(c *gin.Context) {
	postID := c.Param("id")
	if err := updateForumPost(postID, "delete"); err != nil {
		c.JSON(http.StatusBadGateway, apiResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "帖子已删除"})
}

func handleBanUser(c *gin.Context) {
	var req banRequest
	if err := c.ShouldBindJSON(&req); err != nil || (req.Status != "active" && req.Status != "banned") {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Message: "用户状态参数错误"})
		return
	}
	payload, _ := json.Marshal(req)
	resp, err := http.Post(userBaseURL()+"/api/v1/admin/users/"+c.Param("id")+"/ban-sync", "application/json", bytes.NewReader(payload))
	if err != nil {
		c.JSON(http.StatusBadGateway, apiResponse{Success: false, Message: "用户服务不可用"})
		return
	}
	defer resp.Body.Close()
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "用户状态已同步"})
}

func handleStatsOverview(c *gin.Context) {
	forumStats, err := fetchJSON(forumBaseURL() + "/internal/v1/stats/forum-overview")
	if err != nil {
		c.JSON(http.StatusBadGateway, apiResponse{Success: false, Message: err.Error()})
		return
	}
	userStats, err := fetchJSON(userBaseURL() + "/internal/v1/stats/overview")
	if err != nil {
		c.JSON(http.StatusBadGateway, apiResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "获取成功", Data: gin.H{
		"userCount":         userStats["userCount"],
		"todayPostCount":    forumStats["todayPostCount"],
		"pendingAuditCount": forumStats["pendingAuditCount"],
		"postCount":         forumStats["postCount"],
		"boardActivity":     forumStats["boardActivity"],
	}})
}

func reviewAction(c *gin.Context, action string) {
	var req actionRequest
	c.ShouldBindJSON(&req)
	records, err := readAuditRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	recordID := c.Param("id")
	postID := ""
	updated := false
	for i := range records {
		if records[i].ID == recordID {
			records[i].Status = action
			records[i].ReviewerID = req.ReviewerID
			postID = records[i].PostID
			updated = true
		}
	}
	if !updated {
		c.JSON(http.StatusNotFound, apiResponse{Success: false, Message: "审核记录不存在"})
		return
	}
	if err := updateForumPost(postID, action); err != nil {
		c.JSON(http.StatusBadGateway, apiResponse{Success: false, Message: err.Error()})
		return
	}
	if err := writeJSON("audit-records.json", records); err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "审核已处理"})
}

func updateForumPost(postID string, action string) error {
	payload, _ := json.Marshal(gin.H{"action": action})
	resp, err := http.Post(forumBaseURL()+"/internal/v1/posts/"+postID+"/moderate", "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return http.ErrHandlerTimeout
	}
	return nil
}

func fetchJSON(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var envelope struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, err
	}
	return envelope.Data, nil
}

func readConfig() (systemConfig, error) {
	var cfg systemConfig
	return cfg, readJSON("system-config.json", &cfg)
}

func readSensitiveWords() ([]string, error) {
	var words []string
	return words, readJSON("sensitive-words.json", &words)
}

func readAuditRecords() ([]auditRecord, error) {
	var records []auditRecord
	return records, readJSON("audit-records.json", &records)
}

func readPosts() ([]map[string]string, error) {
	var raw []map[string]interface{}
	if err := readJSON("posts.json", &raw); err != nil {
		return nil, err
	}
	items := make([]map[string]string, 0, len(raw))
	for _, item := range raw {
		items = append(items, map[string]string{
			"id":       toString(item["id"]),
			"title":    toString(item["title"]),
			"authorId": toString(item["authorId"]),
			"boardId":  toString(item["boardId"]),
		})
	}
	return items, nil
}

func readJSON(name string, out interface{}) error {
	bytes, err := os.ReadFile(filepath.Join(mockDataDir(), name))
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}

func writeJSON(name string, value interface{}) error {
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(mockDataDir(), name), bytes, 0644)
}

func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}

func mockDataDir() string {
	if custom := os.Getenv("MOCK_DATA_DIR"); custom != "" {
		return custom
	}
	return filepath.Clean(filepath.Join("..", "..", "shared", "mock-data"))
}

func userBaseURL() string {
	return envOr("USER_SERVICE_URL", "http://127.0.0.1:8001")
}

func forumBaseURL() string {
	return envOr("FORUM_SERVICE_URL", "http://127.0.0.1:8002")
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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
