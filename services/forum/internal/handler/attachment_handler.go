package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"ai-forum/forum-service/internal/service"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

// AttachmentHandler handles file attachment operations
type AttachmentHandler struct {
	Service *service.ForumService
}

// NewAttachmentHandler creates a new AttachmentHandler
func NewAttachmentHandler(svc *service.ForumService) *AttachmentHandler {
	return &AttachmentHandler{Service: svc}
}

// UploadAttachment handles file attachment upload
func (h *AttachmentHandler) UploadAttachment(c *gin.Context) {
	userID := c.GetUint("user_id")
	attachmentType := c.PostForm("type") // image, document, link

	if attachmentType == "link" {
		linkURL := c.PostForm("link_url")
		if linkURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "链接地址不能为空"})
			return
		}
		attachment, err := h.Service.UploadAttachment(c.Request.Context(), userID, "link", "link", linkURL, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存链接失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":       attachment.ID,
			"file_type": "link",
			"file_path": linkURL,
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	// Validate file
	if attachmentType == "image" {
		allowedImageTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		}
		mimeType := file.Header.Get("Content-Type")
		if mimeType == "" {
			// Try to detect from file
			src, err := file.Open()
			if err == nil {
				defer src.Close()
				mtype, err := mimetype.DetectReader(src)
				if err == nil {
					mimeType = mtype.String()
				}
			}
		}
		if !allowedImageTypes[mimeType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 JPG/PNG/GIF/WEBP 格式的图片"})
			return
		}
		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过10MB"})
			return
		}
	} else if attachmentType == "document" {
		allowedExtensions := map[string]bool{
			".pdf":  true,
			".doc":  true,
			".docx": true,
			".xlsx": true,
			".txt":  true,
			".md":   true,
		}
		ext := filepath.Ext(file.Filename)
		if !allowedExtensions[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型"})
			return
		}
		if file.Size > 20*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文档大小不能超过20MB"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件类型"})
		return
	}

	// Save file
	now := time.Now()
	uploadDir := fmt.Sprintf("/data/uploads/%d/%02d", now.Year(), now.Month())
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	uniqueName := fmt.Sprintf("%d_%s_%s", userID, now.Format("20060102150405"), filepath.Base(file.Filename))
	destPath := filepath.Join(uploadDir, uniqueName)

	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	filePath := fmt.Sprintf("/uploads/%d/%02d/%s", now.Year(), now.Month(), uniqueName)
	attachment, err := h.Service.UploadAttachment(c.Request.Context(), userID, file.Filename, attachmentType, filePath, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存附件记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        attachment.ID,
		"filename":  attachment.Filename,
		"file_type": attachment.FileType,
		"file_path": attachment.FilePath,
		"file_size": attachment.FileSize,
	})
}

// DownloadAttachment handles file attachment download
func (h *AttachmentHandler) DownloadAttachment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的附件ID"})
		return
	}

	attachment, err := h.Service.GetAttachment(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if attachment.FileType == "link" {
		c.Redirect(http.StatusFound, attachment.FilePath)
		return
	}

	// Serve file from local filesystem
	fullPath := attachment.FilePath
	if !filepath.IsAbs(fullPath) {
		fullPath = "/data" + fullPath
	}

	if attachment.FileType == "image" {
		c.File(fullPath)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachment.Filename))
	c.File(fullPath)
}

// GetPostAttachments returns attachments for a post
func (h *AttachmentHandler) GetPostAttachments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}

	attachments, err := h.Service.GetPostAttachments(c.Request.Context(), uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取附件失败"})
		return
	}
	c.JSON(http.StatusOK, attachments)
}
