package handler

import (
	"net/http"
	"strconv"

	"ai-forum/forum-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ReportAdminHandler struct {
	Service *service.ExtrasService
}

func NewReportAdminHandler(s *service.ExtrasService) *ReportAdminHandler {
	return &ReportAdminHandler{Service: s}
}

func (h *ReportAdminHandler) ListReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.DefaultQuery("status", "pending")
	reports, total, err := h.Service.ListReports(c.Request.Context(), page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": reports, "total": total, "page": page, "limit": limit})
}

func (h *ReportAdminHandler) ResolveReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的举报ID"})
		return
	}
	var req struct {
		Action     string `json:"action" binding:"required"` // resolved | dismissed
		DeletePost bool   `json:"delete_post"`
		AdminNote  string `json:"admin_note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供处理动作"})
		return
	}
	if err := h.Service.ResolveReport(c.Request.Context(), uint(id), req.Action, req.DeletePost, req.AdminNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "举报已处理"})
}
