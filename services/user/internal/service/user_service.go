package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"ai-forum/user-service/internal/model"
	"ai-forum/user-service/pkg/jwt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExpiry  = 30 * time.Minute
	refreshTokenExpiry = 30 * 24 * time.Hour // 30 days
	jwtSecretEnv       = "JWT_SECRET"
)

// UserService handles user business logic
type UserService struct {
	DB  *pgxpool.Pool
	RDB *redis.Client
}

type OAuthTokens struct {
	AccessToken  string
	RefreshToken string
}

// NewUserService creates a new UserService instance
func NewUserService(db *pgxpool.Pool, rdb *redis.Client) *UserService {
	return &UserService{DB: db, RDB: rdb}
}

// resolveJWTReole returns the highest-priority admin role for JWT claims.
func (s *UserService) resolveJWTReole(ctx context.Context, userID uint) string {
	var roleName string
	err := s.DB.QueryRow(ctx, `
		SELECT r.name FROM schema_admin.user_roles ur
		JOIN schema_admin.roles r ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY CASE r.name
			WHEN 'platform_admin' THEN 1
			WHEN 'admin' THEN 2
			ELSE 3
		END
		LIMIT 1
	`, userID).Scan(&roleName)
	if err != nil || roleName == "" {
		return "student"
	}
	if roleName == "user" {
		return "student"
	}
	return roleName
}

// ResolveAppRole returns the frontend-facing role for routing.
func (s *UserService) ResolveAppRole(ctx context.Context, userID uint) string {
	return s.resolveJWTReole(ctx, userID)
}

// Register creates a new user account with invite code validation
func (s *UserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, string, error) {
	// Validate invite code exists and is unused
	var codeStatus string
	var codeID int64
	err := s.DB.QueryRow(ctx,
		"SELECT id, status FROM schema_auth.invite_codes WHERE code = $1",
		req.InvitationCode,
	).Scan(&codeID, &codeStatus)
	if err != nil {
		return nil, "", fmt.Errorf("无效邀请码")
	}
	if codeStatus != "unused" {
		return nil, "", fmt.Errorf("邀请码已被使用或已作废")
	}

	// Check username uniqueness
	var exists bool
	err = s.DB.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM schema_auth.users WHERE username = $1)",
		req.Username,
	).Scan(&exists)
	if err != nil {
		return nil, "", fmt.Errorf("database error: %w", err)
	}
	if exists {
		return nil, "", fmt.Errorf("用户名已被占用")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Start transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert user
	var user model.User
	completed := req.Department != "" && req.Squad != "" && req.Grade != ""
	err = tx.QueryRow(ctx, `
		INSERT INTO schema_auth.users (username, password_hash, nickname, bio, avatar, level, status, department, squad, grade, profile_completed)
		VALUES ($1, $2, $3, '', '', 0, 'active', $4, $5, $6, $7)
		RETURNING id, username, nickname, bio, avatar, level, status, department, squad, grade, profile_completed, created_at, updated_at
	`, req.Username, string(hash), req.Username, req.Department, req.Squad, req.Grade, completed).Scan(
		&user.ID, &user.Username, &user.Nickname, &user.Bio, &user.Avatar,
		&user.Level, &user.Status, &user.Department, &user.Squad, &user.Grade, &user.ProfileCompleted,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Mark invite code as used
	_, err = tx.Exec(ctx,
		"UPDATE schema_auth.invite_codes SET status = 'used', used_by = $1, used_at = NOW() WHERE id = $2",
		user.ID, codeID,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update invite code: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Generate tokens
	secret := os.Getenv(jwtSecretEnv) // Will be overridden by env
	role := s.resolveJWTReole(ctx, user.ID)
	token, err := jwt.GenerateToken(user.ID, user.Username, role, user.Level, secret, accessTokenExpiry)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}
	refreshToken := jwt.GenerateRefreshToken()
	redisKey := fmt.Sprintf("refresh:%d", user.ID)
	s.RDB.Set(ctx, redisKey, refreshToken, refreshTokenExpiry)

	return &user, token, nil
}

// Login authenticates a user and returns JWT tokens
func (s *UserService) Login(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	var user model.User
	err := s.DB.QueryRow(ctx, `
		SELECT id, username, password_hash, nickname, bio, avatar, level, status,
		       department, squad, grade, profile_completed, created_at, updated_at
		FROM schema_auth.users WHERE username = $1
	`, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Nickname, &user.Bio,
		&user.Avatar, &user.Level, &user.Status, &user.Department, &user.Squad, &user.Grade,
		&user.ProfileCompleted, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// Check status
	if user.Status == "banned" {
		return nil, fmt.Errorf("账号已被封禁")
	}

	// Generate tokens
	secret := os.Getenv(jwtSecretEnv)
	role := s.resolveJWTReole(ctx, user.ID)
	accessToken, err := jwt.GenerateToken(user.ID, user.Username, role, user.Level, secret, accessTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken := jwt.GenerateRefreshToken()
	redisKey := fmt.Sprintf("refresh:%d", user.ID)
	s.RDB.Set(ctx, redisKey, refreshToken, refreshTokenExpiry)

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTokenExpiry.Seconds()),
		User:         user.ToResponse(),
	}, nil
}

func (s *UserService) SetOAuthState(ctx context.Context, provider, state, returnTo string, ttl time.Duration) error {
	key := fmt.Sprintf("oauth_state:%s:%s", provider, state)
	return s.RDB.Set(ctx, key, returnTo, ttl).Err()
}

func (s *UserService) ConsumeOAuthState(ctx context.Context, provider, state string) (string, bool, error) {
	key := fmt.Sprintf("oauth_state:%s:%s", provider, state)
	val, err := s.RDB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	_ = s.RDB.Del(ctx, key).Err()
	return val, true, nil
}

// LoginOrCreateOAuthUser logs in with a third-party identity, creating a local user if needed.
// When requireInvite is true, new users are rejected (invite-only mode).
func (s *UserService) LoginOrCreateOAuthUser(ctx context.Context, provider, providerUserID string, rawProfile map[string]any, requireInvite bool) (*model.User, string, *OAuthTokens, error) {
	provider = strings.TrimSpace(provider)
	providerUserID = strings.TrimSpace(providerUserID)
	if provider == "" || providerUserID == "" {
		return nil, "", nil, fmt.Errorf("第三方登录参数错误")
	}

	var userID uint
	err := s.DB.QueryRow(ctx, `
		SELECT user_id FROM schema_auth.oauth_identities
		WHERE provider = $1 AND provider_user_id = $2
	`, provider, providerUserID).Scan(&userID)

	if err != nil && err != pgx.ErrNoRows {
		return nil, "", nil, fmt.Errorf("登录失败，请稍后重试")
	}

	// Existing binding
	if err == nil && userID > 0 {
		user, err := s.GetUserProfile(ctx, userID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("用户不存在")
		}
		if user.Status == "banned" {
			return nil, "", nil, fmt.Errorf("账号已被封禁")
		}
		role := s.resolveJWTReole(ctx, user.ID)
		tokens, err := s.issueTokens(ctx, user.ID, user.Username, role, user.Level)
		if err != nil {
			return nil, "", nil, err
		}
		return user, role, tokens, nil
	}

	if requireInvite {
		return nil, "", nil, fmt.Errorf("当前站点为邀请码模式，请先注册后再绑定 QQ")
	}

	// Create user + binding in a transaction
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("登录失败，请稍后重试")
	}
	defer tx.Rollback(ctx)

	username := buildOAuthUsername(provider, providerUserID)
	// Ensure uniqueness by suffixing random bits if needed
	for i := 0; i < 3; i++ {
		var exists bool
		if err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_auth.users WHERE username = $1)", username).Scan(&exists); err == nil && !exists {
			break
		}
		username = buildOAuthUsername(provider, providerUserID) + "_" + randomShort()
	}

	pwSeed := randomShort() + randomShort()
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwSeed), 12)

	nickname := username
	if v, ok := rawProfile["nickname"].(string); ok && strings.TrimSpace(v) != "" {
		nickname = strings.TrimSpace(v)
	}
	avatar := ""
	for _, k := range []string{"figureurl_qq_2", "figureurl_qq_1", "figureurl"} {
		if v, ok := rawProfile[k].(string); ok && strings.TrimSpace(v) != "" {
			avatar = strings.TrimSpace(v)
			break
		}
	}

	var user model.User
	profileJSON, _ := json.Marshal(rawProfile)
	err = tx.QueryRow(ctx, `
		INSERT INTO schema_auth.users (username, password_hash, nickname, bio, avatar, level, status, department, squad, grade, profile_completed)
		VALUES ($1, $2, $3, '', $4, 0, 'active', '', '', '', false)
		RETURNING id, username, nickname, bio, avatar, level, status,
		          department, squad, grade, profile_completed, created_at, updated_at
	`, username, string(hash), nickname, avatar).Scan(
		&user.ID, &user.Username, &user.Nickname, &user.Bio, &user.Avatar,
		&user.Level, &user.Status, &user.Department, &user.Squad, &user.Grade, &user.ProfileCompleted,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("登录失败，请稍后重试")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO schema_auth.oauth_identities (user_id, provider, provider_user_id, raw_profile)
		VALUES ($1, $2, $3, $4::jsonb)
	`, user.ID, provider, providerUserID, string(profileJSON))
	if err != nil {
		return nil, "", nil, fmt.Errorf("登录失败，请稍后重试")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", nil, fmt.Errorf("登录失败，请稍后重试")
	}

	role := s.resolveJWTReole(ctx, user.ID)
	tokens, err := s.issueTokens(ctx, user.ID, user.Username, role, user.Level)
	if err != nil {
		return nil, "", nil, err
	}
	return &user, role, tokens, nil
}

func (s *UserService) issueTokens(ctx context.Context, userID uint, username, role string, level int) (*OAuthTokens, error) {
	secret := os.Getenv(jwtSecretEnv)
	accessToken, err := jwt.GenerateToken(userID, username, role, level, secret, accessTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	refreshToken := jwt.GenerateRefreshToken()
	redisKey := fmt.Sprintf("refresh:%d", userID)
	s.RDB.Set(ctx, redisKey, refreshToken, refreshTokenExpiry)
	return &OAuthTokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func buildOAuthUsername(provider, providerUserID string) string {
	base := provider + "_" + providerUserID
	if len(base) <= 50 {
		return base
	}
	return base[:50]
}

func randomShort() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// RefreshToken generates a new access token given a valid refresh token
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Find user by refresh token - iterate all refresh tokens
	iter := s.RDB.Scan(ctx, 0, "refresh:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := s.RDB.Get(ctx, key).Result()
		if err == nil && val == refreshToken {
			// Extract user ID from key
			var userID uint
			fmt.Sscanf(key, "refresh:%d", &userID)

			var user model.User
			err := s.DB.QueryRow(ctx,
				"SELECT id, username, level, status FROM schema_auth.users WHERE id = $1",
				userID,
			).Scan(&user.ID, &user.Username, &user.Level, &user.Status)
			if err != nil {
				return "", fmt.Errorf("用户不存在")
			}
			if user.Status == "banned" {
				// Delete refresh token for banned user
				s.RDB.Del(ctx, key)
				return "", fmt.Errorf("账号已被封禁")
			}

			secret := os.Getenv(jwtSecretEnv)
			role := s.resolveJWTReole(ctx, user.ID)
			newToken, err := jwt.GenerateToken(user.ID, user.Username, role, user.Level, secret, accessTokenExpiry)
			if err != nil {
				return "", fmt.Errorf("failed to generate token: %w", err)
			}
			return newToken, nil
		}
	}
	return "", fmt.Errorf("无效或过期令牌")
}

// GetUserProfile retrieves user profile by ID (no password_hash)
func (s *UserService) GetUserProfile(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := s.DB.QueryRow(ctx, `
		SELECT id, username, nickname, bio, avatar, level, status,
		       department, squad, grade, profile_completed, created_at, updated_at
		FROM schema_auth.users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Username, &user.Nickname, &user.Bio, &user.Avatar,
		&user.Level, &user.Status, &user.Department, &user.Squad, &user.Grade,
		&user.ProfileCompleted, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	return &user, nil
}

// UpdateUserProfile updates user profile fields
func (s *UserService) UpdateUserProfile(ctx context.Context, id uint, req *model.UpdateProfileRequest) (*model.User, error) {
	// If username is being updated, check uniqueness
	if req.Username != "" {
		var exists bool
		err := s.DB.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM schema_auth.users WHERE username = $1 AND id != $2)",
			req.Username, id,
		).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("用户名已被占用")
		}
	}

	profileCompleted := req.ProfileCompleted
	var user model.User
	err := s.DB.QueryRow(ctx, `
		UPDATE schema_auth.users
		SET username = COALESCE(NULLIF($1, ''), username),
		    nickname = COALESCE(NULLIF($2, ''), nickname),
		    bio = COALESCE($3, bio),
		    avatar = COALESCE(NULLIF($4, ''), avatar),
		    department = COALESCE(NULLIF($5, ''), department),
		    squad = COALESCE(NULLIF($6, ''), squad),
		    grade = COALESCE(NULLIF($7, ''), grade),
		    profile_completed = COALESCE($8, profile_completed),
		    updated_at = NOW()
		WHERE id = $9
		RETURNING id, username, nickname, bio, avatar, level, status,
		          department, squad, grade, profile_completed, created_at, updated_at
	`, req.Username, req.Nickname, req.Bio, req.Avatar, req.Department, req.Squad, req.Grade,
		profileCompleted, id).Scan(
		&user.ID, &user.Username, &user.Nickname, &user.Bio, &user.Avatar,
		&user.Level, &user.Status, &user.Department, &user.Squad, &user.Grade,
		&user.ProfileCompleted, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	return &user, nil
}

// GetUserStatus returns user status for internal service calls
func (s *UserService) GetUserStatus(ctx context.Context, id uint) (uint, string, int, string, error) {
	var user model.User
	err := s.DB.QueryRow(ctx,
		"SELECT id, username, level, status FROM schema_auth.users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Level, &user.Status)
	if err != nil {
		return 0, "", 0, "", fmt.Errorf("user not found")
	}
	return user.ID, user.Username, user.Level, user.Status, nil
}

// CreateInviteCode generates a single invite code
func (s *UserService) CreateInviteCode(ctx context.Context, createdBy int64) (string, error) {
	code := jwt.GenerateRefreshToken()
	err := s.DB.QueryRow(ctx,
		"INSERT INTO schema_auth.invite_codes (code, created_by) VALUES ($1, $2) RETURNING code",
		code, createdBy,
	).Scan(&code)
	if err != nil {
		return "", fmt.Errorf("failed to create invite code: %w", err)
	}
	return code, nil
}

// CreateInviteCodesBatch generates N invite codes in a transaction
func (s *UserService) CreateInviteCodesBatch(ctx context.Context, count int, createdBy int64) ([]string, error) {
	codes := make([]string, 0, count)
	for i := 0; i < count; i++ {
		code, err := s.CreateInviteCode(ctx, createdBy)
		if err != nil {
			return codes, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}
