package handler

import (
	"net/http"
	"strconv"

	"ai-forum/forum-service/internal/model"
	"ai-forum/forum-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

// ForumHandler holds the service dependency
type ForumHandler struct {
	Service *service.ForumService
}

// NewForumHandler creates a new ForumHandler
func NewForumHandler(svc *service.ForumService) *ForumHandler {
	return &ForumHandler{Service: svc}
}

// ListBoards returns the list of forum boards
func (h *ForumHandler) ListBoards(c *gin.Context) {
	boards, err := h.Service.ListBoards(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取板块列表失败"})
		return
	}
	c.JSON(http.StatusOK, boards)
}

// GetBoard returns a single board by ID
func (h *ForumHandler) GetBoard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的板块ID"})
		return
	}
	board, err := h.Service.GetBoard(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, board)
}

// ListPosts returns a paginated list of posts
func (h *ForumHandler) ListPosts(c *gin.Context) {
	boardIDStr := c.Query("board_id")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	var boardID uint
	if boardIDStr != "" {
		id, err := strconv.ParseUint(boardIDStr, 10, 32)
		if err == nil {
			boardID = uint(id)
		}
	}
	authorIDStr := c.Query("author_id")
	var authorID uint
	if authorIDStr != "" {
		id, err := strconv.ParseUint(authorIDStr, 10, 32)
		if err == nil {
			authorID = uint(id)
		}
	}
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	keyword := c.Query("q")
	sort := c.DefaultQuery("sort", "hot")

	var viewerID uint
	if id, ok := c.Get("user_id"); ok {
		if uid, ok := id.(uint); ok {
			viewerID = uid
		}
	}
	posts, total, err := h.Service.ListPosts(c.Request.Context(), boardID, authorID, page, limit, keyword, sort, viewerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子列表失败"})
		return
	}

	// Batch check follow status for post authors
	if viewerID > 0 && len(posts) > 0 {
		authorIDs := make([]uint, 0, len(posts))
		seen := make(map[uint]bool)
		for _, p := range posts {
			if p.AuthorID > 0 && !seen[p.AuthorID] {
				authorIDs = append(authorIDs, p.AuthorID)
				seen[p.AuthorID] = true
			}
		}
		if followMap, err := h.Service.UserClient.BatchIsFollowing(viewerID, authorIDs); err == nil {
			for i := range posts {
				if posts[i].AuthorID > 0 {
					posts[i].IsFollowing = followMap[posts[i].AuthorID]
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetPost returns a single post by ID
func (h *ForumHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}
	var viewerID uint
	if uid, ok := c.Get("user_id"); ok {
		if v, ok := uid.(uint); ok {
			viewerID = v
		}
	}
	post, err := h.Service.GetPost(c.Request.Context(), uint(id), viewerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Check follow status for post author
	if viewerID > 0 && post.AuthorID > 0 && post.AuthorID != viewerID {
		if followMap, err := h.Service.UserClient.BatchIsFollowing(viewerID, []uint{post.AuthorID}); err == nil {
			post.IsFollowing = followMap[post.AuthorID]
		}
	}

	c.JSON(http.StatusOK, post)
}

// CreatePost creates a new post
func (h *ForumHandler) CreatePost(c *gin.Context) {
	authorID := c.GetUint("user_id")

	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	p := bluemonday.UGCPolicy()
	req.Title = p.Sanitize(req.Title)
	req.Content = p.Sanitize(req.Content)

	post, err := h.Service.CreatePost(c.Request.Context(), authorID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// UpdatePost updates an existing post
func (h *ForumHandler) UpdatePost(c *gin.Context) {
	authorID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	p := bluemonday.UGCPolicy()
	req.Title = p.Sanitize(req.Title)
	req.Content = p.Sanitize(req.Content)

	post, err := h.Service.UpdatePost(c.Request.Context(), authorID, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

// DeletePost deletes a post by ID
func (h *ForumHandler) DeletePost(c *gin.Context) {
	authorID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}

	if err := h.Service.DeletePost(c.Request.Context(), authorID, uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "帖子已删除"})
}
