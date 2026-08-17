package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken parses and validates a JWT token signed by user-service.
// admin-service is a pure consumer of the Session Token Contract: it never
// issues tokens and treats claims as an opaque map.
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

