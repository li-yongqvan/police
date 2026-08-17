package jwt

import (
	"fmt"
	"time"

	"ai-forum/user-service/internal/role"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure. This is the Session Token
// Contract: user_id, username, role, exp, iat (HS256). Do not add fields
// without updating shared/api-contract.md.
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT access token for a user. It refuses to
// issue tokens with a role outside the authorization domain.
func GenerateToken(userID uint, username, roleName, secret string, expiry time.Duration) (string, error) {
	name, ok := role.ValidName(roleName)
	if !ok {
		return "", fmt.Errorf("refusing to issue token with out-of-domain role %q", roleName)
	}

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     string(name),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken parses and validates a JWT token
func ValidateToken(tokenString, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateRefreshToken generates a UUID v4 refresh token string
func GenerateRefreshToken() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(time.Now().UnixNano()>>uint(i*8)) & 0xFF
	}
	// Simple UUID v4 format
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
