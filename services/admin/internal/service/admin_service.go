package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"ai-forum/admin-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AdminService handles admin business logic
type AdminService struct {
	DB             *pgxpool.Pool
	sensitiveWords map[string]string // word -> category
	mu             sync.RWMutex
}

// NewAdminService creates a new AdminService instance
func NewAdminService(db *pgxpool.Pool) *AdminService {
	return &AdminService{DB: db, sensitiveWords: make(map[string]string)}
}

// LoadSensitiveWords loads all sensitive words into the in-memory cache
func (s *AdminService) LoadSensitiveWords(ctx context.Context) error {
	rows, err := s.DB.Query(ctx, "SELECT word, category FROM schema_admin.sensitive_words")
	if err != nil {
		return fmt.Errorf("failed to load sensitive words: %w", err)
	}
	defer rows.Close()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sensitiveWords = make(map[string]string)
	for rows.Next() {
		var word, category string
		if err := rows.Scan(&word, &category); err != nil {
			continue
		}
		s.sensitiveWords[word] = category
	}
	return nil
}

// CheckSensitiveWords checks if text contains any sensitive words
func (s *AdminService) CheckSensitiveWords(text string) (bool, []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matched []string
	for word := range s.sensitiveWords {
		if strings.Contains(text, word) {
			matched = append(matched, word)
		}
	}
	return len(matched) == 0, matched
}

// AddSensitiveWord adds a new sensitive word
func (s *AdminService) AddSensitiveWord(ctx context.Context, word, category string) error {
	_, err := s.DB.Exec(ctx,
		"INSERT INTO schema_admin.sensitive_words (word, category) VALUES ($1, $2) ON CONFLICT (word) DO NOTHING",
		word, category,
	)
	if err != nil {
		return fmt.Errorf("failed to add sensitive word: %w", err)
	}

	s.mu.Lock()
	s.sensitiveWords[word] = category
	s.mu.Unlock()
	return nil
}

// ListSensitiveWords returns all sensitive words
func (s *AdminService) ListSensitiveWords(ctx context.Context) ([]*model.SensitiveWord, error) {
	rows, err := s.DB.Query(ctx, "SELECT id, word, category, created_at FROM schema_admin.sensitive_words ORDER BY word")
	if err != nil {
		return nil, fmt.Errorf("failed to list sensitive words: %w", err)
	}
	defer rows.Close()

	var words []*model.SensitiveWord
	for rows.Next() {
		w := &model.SensitiveWord{}
		if err := rows.Scan(&w.ID, &w.Word, &w.Category, &w.CreatedAt); err != nil {
			continue
		}
		words = append(words, w)
	}
	return words, nil
}

// DeleteSensitiveWord deletes a sensitive word by ID
func (s *AdminService) DeleteSensitiveWord(ctx context.Context, id uint) error {
	_, err := s.DB.Exec(ctx, "DELETE FROM schema_admin.sensitive_words WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete sensitive word: %w", err)
	}

	s.mu.Lock()
	// Reload cache to remove the word
	s.LoadSensitiveWords(ctx)
	s.mu.Unlock()
	return nil
}

// GetConfig retrieves system configuration
func (s *AdminService) GetConfig(key string) (*model.SystemConfig, error) {
	// TODO: implement config retrieval
	return nil, nil
}

// UpdateConfig updates system configuration
func (s *AdminService) UpdateConfig(key, value string) error {
	// TODO: implement config update
	return nil
}

// ListPendingAudit returns posts pending review
func (s *AdminService) ListPendingAudit(page, limit int) ([]*model.AuditRecord, int, error) {
	// TODO: implement pending audit listing
	return nil, 0, nil
}

// GenerateInviteCode generates a single invite code via user-service
func (s *AdminService) GenerateInviteCode(ctx context.Context, userClient *UserClient) (string, error) {
	code, err := userClient.GenerateInviteCode(0)
	if err != nil {
		return "", fmt.Errorf("failed to generate invite code: %w", err)
	}
	return code, nil
}

// GenerateInviteCodesBatch generates multiple invite codes via user-service
func (s *AdminService) GenerateInviteCodesBatch(ctx context.Context, count int, userClient *UserClient) ([]string, error) {
	codes, err := userClient.GenerateInviteCodesBatch(count, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to generate batch invite codes: %w", err)
	}
	return codes, nil
}
