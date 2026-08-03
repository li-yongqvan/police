package handler

import (
	"strings"
	"testing"
	"time"

	"ai-forum/user-service/internal/model"

	"github.com/gin-gonic/gin"
)

func TestAvatarFilenameUsesNanosecondPrecision(t *testing.T) {
	base := time.Date(2026, 7, 20, 14, 0, 0, 0, time.UTC)

	first := avatarFilename(42, "avatar.png", base)
	second := avatarFilename(42, "avatar.png", base.Add(time.Nanosecond))

	if first == second {
		t.Fatalf("avatar filenames should be unique for consecutive uploads: %q", first)
	}
	if !strings.HasPrefix(first, "42_") || !strings.HasSuffix(first, ".png") {
		t.Fatalf("avatar filename should preserve user prefix and extension: %q", first)
	}
}

func TestLoginResponseJSONIncludesResolvedRole(t *testing.T) {
	payload := loginResponseJSON(&model.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    1800,
		User: model.UserResponse{
			ID:       42,
			Username: "demo_admin",
			Nickname: "Demo Admin",
			Level:    2,
			Status:   "active",
		},
	}, "admin")

	if payload["access_token"] != "access-token" {
		t.Fatalf("expected access token to be preserved")
	}
	if payload["token"] != "access-token" {
		t.Fatalf("expected legacy token alias to be preserved")
	}
	if payload["refresh_token"] != "refresh-token" {
		t.Fatalf("expected refresh token to be preserved")
	}
	if payload["expires_in"] != 1800 {
		t.Fatalf("expected expiry to be preserved")
	}

	user, ok := payload["user"].(gin.H)
	if !ok {
		t.Fatalf("expected user payload to be gin.H, got %T", payload["user"])
	}
	if user["role"] != "admin" {
		t.Fatalf("expected resolved role in user payload, got %v", user["role"])
	}
	if user["username"] != "demo_admin" {
		t.Fatalf("expected username to be preserved, got %v", user["username"])
	}
}
