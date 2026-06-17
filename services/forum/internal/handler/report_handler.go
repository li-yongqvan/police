package handler

import (
	"net/http"
	"strconv"

	"ai-forum/forum-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	Service *service.ExtrasService
}

func NewReportHandler(s *service.ExtrasService) *ReportHandler {
	return &ReportHandler{Service: s}
}

func (h *ReportHandler) ReportPost(c *gin.Context) {
	reporterID := c.GetUint("user_id")
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写举报理由"})
		return
	}
	if err := h.Service.ReportPost(c.Request.Context(), reporterID, uint(postID), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "举报提交失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "举报已提交，管理员将尽快处理"})
}
