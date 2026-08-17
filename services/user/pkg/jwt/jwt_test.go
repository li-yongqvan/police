package jwt

import (
	"testing"
	"time"
)

const testSecret = "test-secret"

func TestGenerateTokenEmitsExactSessionTokenContract(t *testing.T) {
	token, err := GenerateToken(7, "alice", "admin", testSecret, 5*time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	claims, err := ValidateToken(token, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	keys := map[string]bool{}
	for k := range claims {
		keys[k] = true
	}
	want := []string{"user_id", "username", "role", "exp", "iat"}
	if len(keys) != len(want) {
		t.Fatalf("claims key set = %v, want exactly %v", keys, want)
	}
	for _, k := range want {
		if !keys[k] {
			t.Fatalf("claim %q missing from %v", k, keys)
		}
	}
	if claims["role"] != "admin" {
		t.Fatalf("role claim = %v, want admin", claims["role"])
	}
	if claims["user_id"] != float64(7) {
		t.Fatalf("user_id claim = %v, want 7", claims["user_id"])
	}
	if claims["username"] != "alice" {
		t.Fatalf("username claim = %v, want alice", claims["username"])
	}
	if _, hasLevel := claims["level"]; hasLevel {
		t.Fatalf("level claim must not be part of the session token contract")
	}
}

func TestGenerateTokenRejectsOutOfDomainRole(t *testing.T) {
	for _, bad := range []string{"moderator", "user", "", "superuser"} {
		if _, err := GenerateToken(7, "alice", bad, testSecret, time.Minute); err == nil {
			t.Fatalf("GenerateToken(%q) succeeded, want refusal", bad)
		}
	}
}

func TestGenerateTokenAcceptsAllDomainRoles(t *testing.T) {
	for _, good := range []string{"student", "admin", "platform_admin"} {
		if _, err := GenerateToken(7, "alice", good, testSecret, time.Minute); err != nil {
			t.Fatalf("GenerateToken(%q) failed: %v", good, err)
		}
	}
}
