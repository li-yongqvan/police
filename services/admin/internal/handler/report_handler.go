package handler

import (
	"net/http"
	"strconv"

	"ai-forum/admin-service/internal/client"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	ForumClient *client.ForumClient
}

func NewReportHandler(fc *client.ForumClient) *ReportHandler {
	return &ReportHandler{ForumClient: fc}
}

func (h *ReportHandler) ListReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.DefaultQuery("status", "pending")
	reports, total, err := h.ForumClient.ListReports(page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": reports, "total": total, "page": page, "limit": limit})
}

func (h *ReportHandler) ResolveReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的举报ID"})
		return
	}
	var req client.ResolveReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.ForumClient.ResolveReport(uint(id), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "举报已处理"})
}
