package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// UserClient handles HTTP communication with user-service
type UserClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewUserClient creates a new UserClient
func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// BatchIsFollowing checks which of the followeeIDs are followed by followerID
func (c *UserClient) BatchIsFollowing(followerID uint, followeeIDs []uint) (map[uint]bool, error) {
	if len(followeeIDs) == 0 {
		return map[uint]bool{}, nil
	}

	type checkReq struct {
		FollowerID  uint   `json:"follower_id"`
		FolloweeIDs []uint `json:"followee_ids"`
	}

	payload := checkReq{FollowerID: followerID, FolloweeIDs: followeeIDs}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/internal/v1/users/following/batch-check", c.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user-service returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Following map[uint]bool `json:"following"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Following, nil
}
