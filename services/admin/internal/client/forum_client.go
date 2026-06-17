package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PendingPost represents a post pending audit
type PendingPost struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	AuthorID     uint      `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	BoardID      uint      `json:"board_id"`
	BoardName    string    `json:"board_name"`
	CreatedAt    time.Time `json:"created_at"`
	MatchedWords []string  `json:"matched_words,omitempty"`
}

// Board represents a forum board (mirrors forum-service model)
type Board struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

// PostSummary represents a post for admin listing
type PostSummary struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	AuthorID   uint      `json:"author_id"`
	AuthorName string    `json:"author_name"`
	BoardID    uint      `json:"board_id"`
	BoardName  string    `json:"board_name"`
	Status     string    `json:"status"`
	IsFeatured bool      `json:"is_featured"`
	IsPinned   bool      `json:"is_pinned"`
	CreatedAt  time.Time `json:"created_at"`
}

// BoardActivity represents activity stats per board
type BoardActivity struct {
	BoardID      uint   `json:"board_id"`
	BoardName    string `json:"board_name"`
	PostCount    int64  `json:"post_count"`
	CommentCount int64  `json:"comment_count"`
}

// DailyCount represents a date-count pair
type DailyCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// ForumOverview contains aggregate forum statistics
type ForumOverview struct {
	TotalPosts    int64 `json:"total_posts"`
	TotalComments int64 `json:"total_comments"`
	TotalLikes    int64 `json:"total_likes"`
	PostsToday    int64 `json:"posts_today"`
	CommentsToday int64 `json:"comments_today"`
}

// ForumClient handles HTTP communication with forum-service for admin operations
type ForumClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewForumClient creates a new ForumClient
func NewForumClient(baseURL string) *ForumClient {
	return &ForumClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// DeletePost calls forum-service to admin-delete a post
func (c *ForumClient) DeletePost(postID uint) error {
	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/internal/v1/posts/%d/admin-delete", c.BaseURL, postID),
		"application/json",
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to call forum-service delete post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("forum-service returned status %d for delete post", resp.StatusCode)
	}
	return nil
}

// SetPostFeatured calls forum-service to set/unset featured flag
func (c *ForumClient) SetPostFeatured(postID uint, featured bool) error {
	payload := map[string]interface{}{"featured": featured}
	return c.postJSON(fmt.Sprintf("%s/internal/v1/posts/%d/admin-featured", c.BaseURL, postID), payload)
}

// SetPostPinned calls forum-service to set/unset pinned flag
func (c *ForumClient) SetPostPinned(postID uint, pinned bool) error {
	payload := map[string]interface{}{"pinned": pinned}
	return c.postJSON(fmt.Sprintf("%s/internal/v1/posts/%d/admin-pinned", c.BaseURL, postID), payload)
}

// ChangePostStatus calls forum-service to change post status (approved/rejected)
func (c *ForumClient) ChangePostStatus(postID uint, status string) error {
	payload := map[string]interface{}{"status": status}
	return c.postJSON(fmt.Sprintf("%s/internal/v1/posts/%d/admin-status", c.BaseURL, postID), payload)
}

// ListPendingPosts calls forum-service to get posts pending review
func (c *ForumClient) ListPendingPosts(page, limit int) ([]PendingPost, int, error) {
	url := fmt.Sprintf("%s/internal/v1/posts/pending?page=%d&limit=%d", c.BaseURL, page, limit)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to call forum-service list pending posts: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("forum-service returned status %d for list pending posts", resp.StatusCode)
	}

	var result struct {
		Posts []PendingPost `json:"posts"`
		Total int           `json:"total"`
		Page  int           `json:"page"`
		Limit int           `json:"limit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Posts, result.Total, nil
}

// BatchDeletePosts calls forum-service to batch delete posts
func (c *ForumClient) BatchDeletePosts(postIDs []uint) error {
	payload := map[string]interface{}{"post_ids": postIDs}
	return c.postJSON(fmt.Sprintf("%s/internal/v1/posts/batch-delete", c.BaseURL), payload)
}

// CreateBoard calls forum-service to create a new board
func (c *ForumClient) CreateBoard(req map[string]interface{}) error {
	return c.postJSON(fmt.Sprintf("%s/internal/v1/boards/admin", c.BaseURL), req)
}

// UpdateBoard calls forum-service to update a board
func (c *ForumClient) UpdateBoard(id uint, req map[string]interface{}) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	url := fmt.Sprintf("%s/internal/v1/boards/admin/%d", c.BaseURL, id)
	req2, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req2.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req2)
	if err != nil {
		return fmt.Errorf("failed to call forum-service update board: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("forum-service returned status %d for update board: %s", resp.StatusCode, string(body))
	}
	return nil
}

// DeleteBoard calls forum-service to delete a board
func (c *ForumClient) DeleteBoard(id uint) error {
	url := fmt.Sprintf("%s/internal/v1/boards/admin/%d", c.BaseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call forum-service delete board: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("forum-service returned status %d for delete board", resp.StatusCode)
	}
	return nil
}

// ListAllBoards calls forum-service to get all boards (including disabled)
func (c *ForumClient) ListAllBoards() ([]Board, error) {
	url := fmt.Sprintf("%s/internal/v1/boards", c.BaseURL)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call forum-service list boards: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forum-service returned status %d for list boards", resp.StatusCode)
	}

	var result struct {
		Boards []Board `json:"boards"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Boards, nil
}

// ListAllPosts calls forum-service to get all posts for admin management
func (c *ForumClient) ListAllPosts(page, limit int) ([]PostSummary, int, error) {
	url := fmt.Sprintf("%s/internal/v1/posts/all?page=%d&limit=%d", c.BaseURL, page, limit)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to call forum-service list posts: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("forum-service returned status %d for list posts", resp.StatusCode)
	}

	var result struct {
		Posts []PostSummary `json:"posts"`
		Total int           `json:"total"`
		Page  int           `json:"page"`
		Limit int           `json:"limit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Posts, result.Total, nil
}

// GetForumOverview calls forum-service internal stats overview
func (c *ForumClient) GetForumOverview() (*ForumOverview, error) {
	url := fmt.Sprintf("%s/internal/v1/stats/overview", c.BaseURL)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call forum-service stats: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forum-service returned status %d for stats", resp.StatusCode)
	}
	var result struct {
		Data ForumOverview `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &result.Data, nil
}

// GetDailyPosts calls forum-service internal daily posts endpoint
func (c *ForumClient) GetDailyPosts(days int) ([]DailyCount, error) {
	url := fmt.Sprintf("%s/internal/v1/stats/daily-posts?days=%d", c.BaseURL, days)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call forum-service daily posts: %w", err)
	}
	defer resp.Body.Close()
	var result struct {
		Data []DailyCount `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data, nil
}

// ReportItem is a user report on a post.
type ReportItem struct {
	ID           uint   `json:"id"`
	PostID       uint   `json:"post_id"`
	PostTitle    string `json:"post_title"`
	ReporterID   uint   `json:"reporter_id"`
	ReporterName string `json:"reporter_name"`
	Reason       string `json:"reason"`
	Status       string `json:"status"`
	AdminNote    string `json:"admin_note,omitempty"`
	CreatedAt    string `json:"created_at"`
}

// ResolveReportRequest is the body for resolving a report.
type ResolveReportRequest struct {
	Action     string `json:"action"`
	DeletePost bool   `json:"delete_post"`
	AdminNote  string `json:"admin_note"`
}

// ListReports returns paginated post reports from forum-service.
func (c *ForumClient) ListReports(page, limit int, status string) ([]ReportItem, int, error) {
	url := fmt.Sprintf("%s/internal/v1/reports?page=%d&limit=%d&status=%s", c.BaseURL, page, limit, status)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list reports: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("forum-service returned status %d for list reports", resp.StatusCode)
	}
	var result struct {
		Reports []ReportItem `json:"reports"`
		Total   int          `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode reports: %w", err)
	}
	return result.Reports, result.Total, nil
}

// ResolveReport marks a report resolved or dismissed.
func (c *ForumClient) ResolveReport(reportID uint, req ResolveReportRequest) error {
	return c.postJSON(fmt.Sprintf("%s/internal/v1/reports/%d/resolve", c.BaseURL, reportID), req)
}

// GetBoardActivity calls forum-service internal board activity endpoint
func (c *ForumClient) GetBoardActivity() ([]BoardActivity, error) {
	url := fmt.Sprintf("%s/internal/v1/stats/board-activity", c.BaseURL)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call forum-service board activity: %w", err)
	}
	defer resp.Body.Close()
	var result struct {
		Data []BoardActivity `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data, nil
}

func (c *ForumClient) postJSON(url string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to call forum-service: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("forum-service returned status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
