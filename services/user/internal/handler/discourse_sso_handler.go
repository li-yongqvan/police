package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"ai-forum/user-service/internal/service"
	"ai-forum/user-service/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// DiscourseSSOHandler handles Discourse Connect (SSO) protocol requests.
type DiscourseSSOHandler struct {
	Service *service.UserService
}

// NewDiscourseSSOHandler creates a new DiscourseSSOHandler.
func NewDiscourseSSOHandler(svc *service.UserService) *DiscourseSSOHandler {
	return &DiscourseSSOHandler{Service: svc}
}

// HandleSSO implements the full Discourse Connect SSO flow:
// 1. Validate HMAC-SHA256 signature of incoming payload
// 2. Authenticate user via JWT cookie
// 3. Build return payload with user attributes
// 4. Sign and redirect back to Discourse
func (h *DiscourseSSOHandler) HandleSSO(c *gin.Context) {
	rawPayload := c.Query("sso")
	sig := c.Query("sig")

	if rawPayload == "" || sig == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing sso or sig parameter"})
		return
	}

	secret := os.Getenv("DISCOURSE_CONNECT_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SSO not configured"})
		return
	}

	// Verify incoming HMAC signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(rawPayload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	// Decode incoming payload: Base64 -> URL-decode -> parse nonce and return_sso_url
	payloadBytes, err := base64.StdEncoding.DecodeString(rawPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sso payload encoding"})
		return
	}

	payloadStr, err := url.QueryUnescape(string(payloadBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sso payload"})
		return
	}

	incoming, err := url.ParseQuery(payloadStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sso payload"})
		return
	}

	nonce := incoming.Get("nonce")
	returnSSOURL := incoming.Get("return_sso_url")
	if nonce == "" || returnSSOURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing nonce or return_sso_url"})
		return
	}

	// Authenticate user via JWT cookie
	cookie, err := c.Request.Cookie("ai-forum-token")
	if err != nil {
		// No valid session cookie, redirect to login
		loginURL := "/"
		c.Redirect(http.StatusFound, loginURL)
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}

	claims, err := jwt.ValidateToken(cookie.Value, jwtSecret)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.Redirect(http.StatusFound, "/")
		return
	}
	userID := uint(userIDFloat)

	// Query user profile
	user, err := h.Service.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user.Status == "banned" {
		c.JSON(http.StatusForbidden, gin.H{"error": "account banned"})
		return
	}

	// Resolve role for Discourse admin/moderator mapping
	role := h.Service.ResolveAppRole(c.Request.Context(), userID)

	// Build return payload
	returnPayload := make(url.Values)
	returnPayload.Set("nonce", nonce)
	returnPayload.Set("external_id", strconv.FormatUint(uint64(user.ID), 10))
	returnPayload.Set("email", user.Username+"@ai-forum.local")

	// Name: prefer Nickname, fall back to Username
	name := user.Nickname
	if name == "" {
		name = user.Username
	}
	returnPayload.Set("name", name)

	// Username for Discourse
	returnPayload.Set("username", user.Username)

	// Bio: truncate to 3000 runes if present
	if user.Bio != "" {
		bio := []rune(user.Bio)
		if len(bio) > 3000 {
			bio = bio[:3000]
		}
		returnPayload.Set("bio", string(bio))
	}

	// Avatar URL: construct from PLATFORM_BASE_URL + relative avatar path
	if user.Avatar != "" {
		avatarURL := strings.TrimSpace(user.Avatar)
		if !strings.HasPrefix(avatarURL, "http://") && !strings.HasPrefix(avatarURL, "https://") {
			platformBaseURL := os.Getenv("PLATFORM_BASE_URL")
			if platformBaseURL == "" {
				platformBaseURL = "http://122.51.233.225:8888"
			}
			baseURL := strings.TrimRight(platformBaseURL, "/")
			if strings.HasPrefix(avatarURL, "/uploads/") && !strings.HasSuffix(baseURL, "/user-api") {
				avatarURL = "/user-api" + avatarURL
			}
			avatarURL = baseURL + "/" + strings.TrimLeft(avatarURL, "/")
		}
		returnPayload.Set("avatar_url", avatarURL)
	}

	// Discourse admin/mod mapping
	if role == "platform_admin" {
		returnPayload.Set("admin", "true")
	}
	if role == "admin" {
		returnPayload.Set("moderator", "true")
	}

	// Encode return payload: URL-encode -> Base64
	returnQuery := returnPayload.Encode()
	returnEncoded := base64.StdEncoding.EncodeToString([]byte(returnQuery))

	// Sign with HMAC-SHA256
	mac.Reset()
	mac.Write([]byte(returnEncoded))
	returnSig := hex.EncodeToString(mac.Sum(nil))

	// Build redirect URL
	targetURL := fmt.Sprintf("%s?sso=%s&sig=%s", returnSSOURL, url.QueryEscape(returnEncoded), returnSig)
	c.Redirect(http.StatusFound, targetURL)
}
