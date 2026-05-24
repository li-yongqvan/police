package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UserClient handles HTTP communication with user-service for invite code generation
type UserClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewUserClient creates a new UserClient
func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GenerateInviteCode generates a single invite code via user-service internal API
func (c *UserClient) GenerateInviteCode(createdBy int64) (string, error) {
	payload := map[string]interface{}{"created_by": createdBy}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/internal/v1/invite-codes", c.BaseURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	var result struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Code, nil
}

// GenerateInviteCodesBatch generates N invite codes via user-service internal API
func (c *UserClient) GenerateInviteCodesBatch(count int, createdBy int64) ([]string, error) {
	payload := map[string]interface{}{"count": count, "created_by": createdBy}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/internal/v1/invite-codes/batch", c.BaseURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	var result struct {
		Codes []string `json:"codes"`
		Count int      `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Codes, nil
}

// BanUser calls user-service to ban a user
func (c *UserClient) BanUser(userID uint, reason string) error {
	payload := map[string]interface{}{"reason": reason}
	return c.postJSON(fmt.Sprintf("%s/internal/v1/users/%d/ban", c.BaseURL, userID), payload)
}

// ListUsers calls user-service to get paginated user list
func (c *UserClient) ListUsers(page, limit int, status string) ([]map[string]interface{}, int, error) {
	url := fmt.Sprintf("%s/internal/v1/users?page=%d&limit=%d", c.BaseURL, page, limit)
	if status != "" {
		url += "&status=" + status
	}
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to call user-service list users: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("user-service returned status %d for list users", resp.StatusCode)
	}

	var result struct {
		Users []map[string]interface{} `json:"users"`
		Total int                      `json:"total"`
		Page  int                      `json:"page"`
		Limit int                      `json:"limit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Users, result.Total, nil
}

// UpdateUserLevel calls user-service to update a user's level
func (c *UserClient) UpdateUserLevel(userID uint, newLevel int) error {
	payload := map[string]interface{}{"level": newLevel}
	return c.putJSON(fmt.Sprintf("%s/internal/v1/users/%d/level", c.BaseURL, userID), payload)
}

// GetUserLogs calls user-service to get operation logs for a user
func (c *UserClient) GetUserLogs(userID uint, page, limit int) ([]map[string]interface{}, int, error) {
	url := fmt.Sprintf("%s/internal/v1/users/%d/logs?page=%d&limit=%d", c.BaseURL, userID, page, limit)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to call user-service get user logs: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("user-service returned status %d for get user logs", resp.StatusCode)
	}

	var result struct {
		Logs  []map[string]interface{} `json:"logs"`
		Total int                      `json:"total"`
		Page  int                      `json:"page"`
		Limit int                      `json:"limit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Logs, result.Total, nil
}

// GetInviteCodeStatus calls user-service to get invite code details
func (c *UserClient) GetInviteCodeStatus(code string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/internal/v1/invite-codes/%s/status", c.BaseURL, code)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service get invite code status: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status %d for get invite code status", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result, nil
}

// ListInviteCodes calls user-service to get paginated invite codes
func (c *UserClient) ListInviteCodes(page, limit int) ([]map[string]interface{}, int, error) {
	url := fmt.Sprintf("%s/internal/v1/invite-codes?page=%d&limit=%d", c.BaseURL, page, limit)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to call user-service list invite codes: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("user-service returned status %d for list invite codes", resp.StatusCode)
	}

	var result struct {
		Codes []map[string]interface{} `json:"codes"`
		Total int                      `json:"total"`
		Page  int                      `json:"page"`
		Limit int                      `json:"limit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Codes, result.Total, nil
}

// VoidInviteCode calls user-service to void an invite code
func (c *UserClient) VoidInviteCode(code string) error {
	return c.putJSON(fmt.Sprintf("%s/internal/v1/invite-codes/%s/void", c.BaseURL, code), nil)
}

// GetUserOverview calls user-service internal stats overview
func (c *UserClient) GetUserOverview() (*UserOverview, error) {
	url := fmt.Sprintf("%s/internal/v1/stats/overview", c.BaseURL)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service stats: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status %d for stats", resp.StatusCode)
	}
	var result struct {
		Data UserOverview `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &result.Data, nil
}

// GetDailyUsers calls user-service internal daily users endpoint
func (c *UserClient) GetDailyUsers(days int) ([]DailyUserCount, error) {
	url := fmt.Sprintf("%s/internal/v1/stats/daily-users?days=%d", c.BaseURL, days)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service daily users: %w", err)
	}
	defer resp.Body.Close()
	var result struct {
		Data []DailyUserCount `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data, nil
}

// GetLevelDistribution calls user-service internal level distribution endpoint
func (c *UserClient) GetLevelDistribution() ([]LevelDistribution, error) {
	url := fmt.Sprintf("%s/internal/v1/stats/level-distribution", c.BaseURL)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service level distribution: %w", err)
	}
	defer resp.Body.Close()
	var result struct {
		Data []LevelDistribution `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Data, nil
}

func (c *UserClient) postJSON(url string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}
	return nil
}

func (c *UserClient) putJSON(url string, payload interface{}) error {
	var bodyReader *bytes.Buffer
	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		bodyReader = bytes.NewBuffer(body)
	} else {
		bodyReader = bytes.NewBuffer([]byte("{}"))
	}
	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}
	return nil
}
