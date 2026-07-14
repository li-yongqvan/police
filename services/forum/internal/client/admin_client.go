package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AdminClient handles HTTP communication with admin-service
type AdminClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAdminClient creates a new AdminClient
func NewAdminClient(baseURL string) *AdminClient {
	return &AdminClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CheckSensitiveWords sends text to admin-service for sensitive word check
func (c *AdminClient) CheckSensitiveWords(text string) (bool, []string, error) {
	payload := map[string]string{"text": text}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return false, nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/internal/v1/moderation/check", c.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		return false, nil, fmt.Errorf("failed to call admin-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, nil, fmt.Errorf("admin-service returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Clean        bool     `json:"clean"`
		MatchedWords []string `json:"matched_words"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Clean, result.MatchedWords, nil
}
