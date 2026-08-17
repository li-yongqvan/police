package role

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPFetcher fetches the Authoritative Role Name from the Role Authority
// internal endpoint (GET /internal/v1/users/:id/role).
type HTTPFetcher struct {
	baseURL string
	client  *http.Client
}

func NewHTTPFetcher(baseURL string, timeout time.Duration) *HTTPFetcher {
	return &HTTPFetcher{baseURL: baseURL, client: &http.Client{Timeout: timeout}}
}

type roleResponse struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

func (f *HTTPFetcher) Fetch(ctx context.Context, userID uint) (Name, error) {
	url := fmt.Sprintf("%s/internal/v1/users/%d/role", f.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build role request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("role authority request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("role authority returned status %d", resp.StatusCode)
	}

	var body roleResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("failed to decode role response: %w", err)
	}
	if body.UserID != userID {
		return "", fmt.Errorf("role response user_id mismatch: got %d, want %d", body.UserID, userID)
	}

	name, ok := ValidName(body.Role)
	if !ok {
		return "", fmt.Errorf("role authority returned out-of-domain role %q", body.Role)
	}
	return name, nil
}
