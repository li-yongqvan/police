package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"ai-forum/user-service/internal/model"
	"ai-forum/user-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

// UserHandler holds the service dependency
type UserHandler struct {
	Service *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{Service: svc}
}

// Register handles user registration with invite code
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	user, _, err := h.Service.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user":    user.ToResponse(),
	})
}

// Login handles user authentication
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	resp, err := h.Service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken handles access token refresh
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	newToken, err := h.Service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newToken,
		"expires_in":   1800,
	})
}

// GetProfile returns user profile information
func (h *UserHandler) GetProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := h.Service.GetUserProfile(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	role := h.Service.ResolveAppRole(c.Request.Context(), uint(id))
	c.JSON(http.StatusOK, toUserJSON(user.ToResponse(), role))
}

// UpdateProfile updates user profile information
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// Verify ownership
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)
	if userID != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改他人资料"})
		return
	}

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	p := bluemonday.UGCPolicy()
	req.Nickname = p.Sanitize(req.Nickname)
	req.Bio = p.Sanitize(req.Bio)

	user, err := h.Service.UpdateUserProfile(c.Request.Context(), uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := h.Service.ResolveAppRole(c.Request.Context(), uint(id))
	c.JSON(http.StatusOK, toUserJSON(user.ToResponse(), role))
}

// UploadAvatar handles avatar image upload
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// Verify ownership
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)
	if userID != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改他人头像"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择头像文件"})
		return
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 JPG/PNG/GIF/WEBP 格式的图片"})
		return
	}

	// Check file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过5MB"})
		return
	}

	// Save to local filesystem
	ext := filepath.Ext(file.Filename)
	filename := strconv.FormatUint(id, 10) + "_" + time.Now().Format("20060102150405") + ext
	uploadDir := "/data/uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Failed to create avatar directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	destPath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	avatarPath := "/uploads/avatars/" + filename

	// Update database
	req := &model.UpdateProfileRequest{Avatar: avatarPath}
	user, err := h.Service.UpdateUserProfile(c.Request.Context(), uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "头像更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "头像上传成功",
		"avatar":  user.Avatar,
	})
}

var demoRoleUsernames = map[string]string{
	"student":         "demo_student",
	"admin":           "demo_admin",
	"platform_admin":  "demo_platform_admin",
}

var demoRoleFromUsername = map[string]string{
	"demo_student":         "student",
	"demo_admin":           "admin",
	"demo_platform_admin":  "platform_admin",
}

// DemoLogin supports MVP one-click role login (password: demo123456).
func (h *UserHandler) DemoLogin(c *gin.Context) {
	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "登录参数不完整"})
		return
	}

	username, ok := demoRoleUsernames[req.Role]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未知演示角色"})
		return
	}

	resp, err := h.Service.Login(c.Request.Context(), username, "demo123456")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	role := h.Service.ResolveAppRole(c.Request.Context(), resp.User.ID)
	if role == "" {
		role = req.Role
	}
	c.JSON(http.StatusOK, gin.H{
		"token":         resp.AccessToken,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"user":          toUserJSON(resp.User, role),
	})
}

// Me returns the current authenticated user profile.
func (h *UserHandler) Me(c *gin.Context) {
	userIDAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := userIDAny.(uint)

	user, err := h.Service.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	role := h.Service.ResolveAppRole(c.Request.Context(), userID)
	c.JSON(http.StatusOK, toUserJSON(user.ToResponse(), role))
}

func toUserJSON(u model.UserResponse, role string) gin.H {
	name := u.Nickname
	if u.Squad != "" && name != "" {
		name = u.Squad + "·" + name
	}
	return gin.H{
		"id":                strconv.FormatUint(uint64(u.ID), 10),
		"name":              name,
		"username":          u.Username,
		"avatar":            u.Avatar,
		"role":              role,
		"department":        u.Department,
		"squad":             u.Squad,
		"grade":             u.Grade,
		"profile_completed": u.ProfileCompleted,
		"bio":               u.Bio,
		"status":            u.Status,
		"level":             u.Level,
	}
}
