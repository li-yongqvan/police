package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ai-forum/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *service.UserService
}

func NewAuthHandler(svc *service.UserService) *AuthHandler {
	return &AuthHandler{Service: svc}
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	user, err := h.Service.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[logout-sync] get user profile failed for user_id=%d: %v", userID, err)
	} else {
		go h.syncLogoutToDiscourse(user.Username)
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) syncLogoutToDiscourse(username string) {
	discourseURL := os.Getenv("DISCOURSE_URL")
	apiKey := os.Getenv("DISCOURSE_API_KEY")
	if discourseURL == "" || apiKey == "" || username == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	discourseUserID, err := h.lookupDiscourseUserByUsername(ctx, discourseURL, username)
	if err != nil {
		log.Printf("[logout-sync] lookup Discourse user by username=%s failed: %v", username, err)
		return
	}

	if err := h.logoutDiscourseUser(ctx, discourseURL, apiKey, discourseUserID); err != nil {
		log.Printf("[logout-sync] Discourse logout user_id=%d failed: %v", discourseUserID, err)
	}
}

func (h *AuthHandler) lookupDiscourseUserByUsername(ctx context.Context, baseURL, username string) (int, error) {
	url := fmt.Sprintf("%s/u/%s.json", baseURL, username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, fmt.Errorf("user %s not found on Discourse", username)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		User struct {
			ID int `json:"id"`
		} `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decode response: %w", err)
	}

	if result.User.ID == 0 {
		return 0, fmt.Errorf("user not found in Discourse response")
	}

	return result.User.ID, nil
}

func (h *AuthHandler) logoutDiscourseUser(ctx context.Context, baseURL, apiKey string, discourseUserID int) error {
	url := fmt.Sprintf("%s/admin/users/%d/log_out", baseURL, discourseUserID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Api-Username", "system")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return nil
}
