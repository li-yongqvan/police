package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type board struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type attachment struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type post struct {
	ID           string       `json:"id"`
	BoardID      string       `json:"boardId"`
	AuthorID     string       `json:"authorId"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	Tags         []string     `json:"tags"`
	Attachments  []attachment `json:"attachments"`
	Status       string       `json:"status"`
	IsFeatured   bool         `json:"isFeatured"`
	LikeCount    int          `json:"likeCount"`
	CommentCount int          `json:"commentCount"`
	CreatedAt    string       `json:"createdAt"`
}

type comment struct {
	ID        string `json:"id"`
	PostID    string `json:"postId"`
	AuthorID  string `json:"authorId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type systemConfig struct {
	PostingEnabled bool            `json:"postingEnabled"`
	BoardSwitches  map[string]bool `json:"boardSwitches"`
	ModerationMode string          `json:"moderationMode"`
}

type postRequest struct {
	BoardID     string       `json:"boardId"`
	AuthorID    string       `json:"authorId"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	Tags        []string     `json:"tags"`
	Attachments []attachment `json:"attachments"`
}

type commentRequest struct {
	AuthorID string `json:"authorId"`
	Content  string `json:"content"`
}

type moderateRequest struct {
	Action string `json:"action"`
}

type moderationCheckResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type forumStats struct {
	TodayPostCount    int     `json:"todayPostCount"`
	PendingAuditCount int     `json:"pendingAuditCount"`
	BoardActivity     []gin.H `json:"boardActivity"`
	PostCount         int     `json:"postCount"`
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

	router.GET("/api/v1/boards", handleBoards)
	router.GET("/api/v1/posts", handlePosts)
	router.GET("/api/v1/posts/:id", handlePostDetail)
	router.POST("/api/v1/posts", handleCreatePost)
	router.POST("/api/v1/posts/:id/comments", handleCreateComment)
	router.POST("/api/v1/posts/:id/like", handleLikePost)
	router.POST("/api/v1/posts/:id/feature", handleFeaturePost)
	router.POST("/internal/v1/posts/:id/moderate", handleModeratePost)
	router.GET("/internal/v1/stats/forum-overview", handleStats)

	port := envOr("FORUM_SERVICE_PORT", "8002")
	router.Run(":" + port)
}

func handleBoards(c *gin.Context) {
	boards, err := readBoards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	config, _ := readConfig()
	includeDisabled := c.Query("includeDisabled") == "1"
	visible := make([]board, 0, len(boards))
	for _, item := range boards {
		if includeDisabled || config.BoardSwitches[item.ID] {
			visible = append(visible, item)
		}
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "获取成功", Data: visible})
}

func handlePosts(c *gin.Context) {
	posts, err := readPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	boardID := c.Query("boardId")
	includePending := c.Query("includePending") == "1"
	filtered := make([]post, 0, len(posts))
	for _, item := range posts {
		if boardID != "" && item.BoardID != boardID {
			continue
		}
		if !includePending && item.Status != "published" {
			continue
		}
		if includePending && !slices.Contains([]string{"published", "pending", "rejected"}, item.Status) {
			continue
		}
		filtered = append(filtered, item)
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "获取成功", Data: filtered})
}

func handlePostDetail(c *gin.Context) {
	posts, err := readPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	postID := c.Param("id")
	includeHidden := c.Query("includeHidden") == "1"
	for _, item := range posts {
		if item.ID == postID {
			if item.Status != "published" && !includeHidden {
				c.JSON(http.StatusNotFound, apiResponse{Success: false, Message: "帖子暂不可见"})
				return
			}
			comments, _ := readComments()
			related := []comment{}
			for _, line := range comments {
				if line.PostID == postID {
					related = append(related, line)
				}
			}
			c.JSON(http.StatusOK, apiResponse{Success: true, Message: "获取成功", Data: gin.H{
				"post":     item,
				"comments": related,
			}})
			return
		}
	}
	c.JSON(http.StatusNotFound, apiResponse{Success: false, Message: "帖子不存在"})
}

func handleCreatePost(c *gin.Context) {
	var req postRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Message: "发帖参数不完整"})
		return
	}

	config, err := readConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	if !config.PostingEnabled {
		c.JSON(http.StatusForbidden, apiResponse{Success: false, Message: "当前已关闭发帖"})
		return
	}
	if !config.BoardSwitches[req.BoardID] {
		c.JSON(http.StatusForbidden, apiResponse{Success: false, Message: "该板块当前未开放"})
		return
	}

	allowed, err := userAllowed(req.AuthorID)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiResponse{Success: false, Message: "用户状态校验失败"})
		return
	}
	if !allowed {
		c.JSON(http.StatusForbidden, apiResponse{Success: false, Message: "该用户已被封禁，无法发帖"})
		return
	}

	posts, err := readPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}

	newPost := post{
		ID:           "post-" + time.Now().Format("20060102150405"),
		BoardID:      req.BoardID,
		AuthorID:     req.AuthorID,
		Title:        strings.TrimSpace(req.Title),
		Content:      strings.TrimSpace(req.Content),
		Tags:         req.Tags,
		Attachments:  decorateAttachments(req.Attachments),
		Status:       "published",
		IsFeatured:   false,
		LikeCount:    0,
		CommentCount: 0,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	checkResult, err := moderationCheck(newPost)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiResponse{Success: false, Message: "审核服务调用失败"})
		return
	}
	newPost.Status = checkResult.Status
	posts = append([]post{newPost}, posts...)
	if err := writePosts(posts); err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}

	message := "帖子已发布"
	if newPost.Status == "pending" {
		message = "帖子已提交，等待审核"
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: message, Data: newPost})
}

func handleCreateComment(c *gin.Context) {
	var req commentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Message: "评论参数不完整"})
		return
	}

	posts, err := readPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	postID := c.Param("id")
	found := false
	for i := range posts {
		if posts[i].ID == postID && posts[i].Status == "published" {
			posts[i].CommentCount++
			found = true
		}
	}
	if !found {
		c.JSON(http.StatusNotFound, apiResponse{Success: false, Message: "帖子不存在或不可评论"})
		return
	}

	comments, err := readComments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	newComment := comment{
		ID:        "comment-" + time.Now().Format("20060102150405"),
		PostID:    postID,
		AuthorID:  req.AuthorID,
		Content:   strings.TrimSpace(req.Content),
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	comments = append(comments, newComment)
	if err := writeComments(comments); err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	if err := writePosts(posts); err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "评论已发布", Data: newComment})
}

func handleLikePost(c *gin.Context) {
	posts, err := readPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	postID := c.Param("id")
	for i := range posts {
		if posts[i].ID == postID {
			posts[i].LikeCount++
			if err := writePosts(posts); err != nil {
				c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
				return
			}
			c.JSON(http.StatusOK, apiResponse{Success: true, Message: "点赞成功", Data: posts[i]})
			return
		}
	}
	c.JSON(http.StatusNotFound, apiResponse{Success: false, Message: "帖子不存在"})
}

func handleFeaturePost(c *gin.Context) {
	posts, err := readPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	postID := c.Param("id")
	for i := range posts {
		if posts[i].ID == postID {
			posts[i].IsFeatured = !posts[i].IsFeatured
			if err := writePosts(posts); err != nil {
				c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
				return
			}
			c.JSON(http.StatusOK, apiResponse{Success: true, Message: "精华状态已切换", Data: posts[i]})
			return
		}
	}
	c.JSON(http.StatusNotFound, apiResponse{Success: false, Message: "帖子不存在"})
}

func handleModeratePost(c *gin.Context) {
	var req moderateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Message: "审核动作不正确"})
		return
	}

	posts, err := readPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	postID := c.Param("id")
	for i := range posts {
		if posts[i].ID == postID {
			switch req.Action {
			case "approve":
				posts[i].Status = "published"
			case "reject":
				posts[i].Status = "rejected"
			case "delete":
				posts[i].Status = "deleted"
			default:
				c.JSON(http.StatusBadRequest, apiResponse{Success: false, Message: "未知审核动作"})
				return
			}
			if err := writePosts(posts); err != nil {
				c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
				return
			}
			c.JSON(http.StatusOK, apiResponse{Success: true, Message: "帖子状态已更新", Data: posts[i]})
			return
		}
	}
	c.JSON(http.StatusNotFound, apiResponse{Success: false, Message: "帖子不存在"})
}

func handleStats(c *gin.Context) {
	posts, err := readPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Message: err.Error()})
		return
	}
	boards, _ := readBoards()
	start := time.Now().Format("2006-01-02")
	boardActivity := make([]gin.H, 0, len(boards))
	publishedCount := 0
	todayCount := 0
	pendingCount := 0
	boardCounts := map[string]int{}
	for _, item := range posts {
		if item.Status == "published" {
			publishedCount++
			boardCounts[item.BoardID]++
		}
		if item.Status == "pending" {
			pendingCount++
		}
		if strings.HasPrefix(item.CreatedAt, start) {
			todayCount++
		}
	}
	for _, b := range boards {
		boardActivity = append(boardActivity, gin.H{
			"boardId": b.ID,
			"name":    b.Name,
			"count":   boardCounts[b.ID],
		})
	}
	c.JSON(http.StatusOK, apiResponse{Success: true, Message: "获取成功", Data: forumStats{
		TodayPostCount:    todayCount,
		PendingAuditCount: pendingCount,
		BoardActivity:     boardActivity,
		PostCount:         publishedCount,
	}})
}

func moderationCheck(item post) (moderationCheckResponse, error) {
	payload, _ := json.Marshal(item)
	resp, err := http.Post(adminBaseURL()+"/internal/v1/moderation/check-post", "application/json", bytes.NewReader(payload))
	if err != nil {
		return moderationCheckResponse{}, err
	}
	defer resp.Body.Close()
	var envelope struct {
		Success bool                    `json:"success"`
		Data    moderationCheckResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return moderationCheckResponse{}, err
	}
	if !envelope.Success {
		return moderationCheckResponse{}, errors.New("moderation check failed")
	}
	return envelope.Data, nil
}

func userAllowed(userID string) (bool, error) {
	resp, err := http.Get(userBaseURL() + "/internal/v1/users/" + userID + "/status")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var envelope struct {
		Success bool `json:"success"`
		Data    struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return false, err
	}
	return envelope.Data.Status != "banned", nil
}

func decorateAttachments(items []attachment) []attachment {
	decorated := make([]attachment, 0, len(items))
	for index, item := range items {
		item.ID = "att-" + time.Now().Format("150405") + "-" + string(rune('a'+index))
		if item.URL == "" {
			item.URL = "https://example.com/resource/" + item.Name
		}
		decorated = append(decorated, item)
	}
	return decorated
}

func readBoards() ([]board, error) {
	var items []board
	return items, readJSON("boards.json", &items)
}

func readPosts() ([]post, error) {
	var items []post
	return items, readJSON("posts.json", &items)
}

func writePosts(items []post) error {
	return writeJSON("posts.json", items)
}

func readComments() ([]comment, error) {
	var items []comment
	return items, readJSON("comments.json", &items)
}

func writeComments(items []comment) error {
	return writeJSON("comments.json", items)
}

func readConfig() (systemConfig, error) {
	var cfg systemConfig
	return cfg, readJSON("system-config.json", &cfg)
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

func mockDataDir() string {
	if custom := os.Getenv("MOCK_DATA_DIR"); custom != "" {
		return custom
	}
	return filepath.Clean(filepath.Join("..", "..", "shared", "mock-data"))
}

func adminBaseURL() string {
	return envOr("ADMIN_SERVICE_URL", "http://127.0.0.1:8003")
}

func userBaseURL() string {
	return envOr("USER_SERVICE_URL", "http://127.0.0.1:8001")
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
