package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testJWTSecret = "arch-deep-4-test-secret"

func signMapClaimsToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

func newAuthProbeRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/probe", func(c *gin.Context) {
		uid, ok := c.Get("user_id")
		if !ok {
			c.String(http.StatusInternalServerError, "user_id not set")
			return
		}
		c.String(http.StatusOK, "uid=%d", uid.(uint))
	})
	return r
}

func TestAuthMiddlewareAllowsAdminAndSetsUserID(t *testing.T) {
	t.Setenv("JWT_SECRET", testJWTSecret)
	r := newAuthProbeRouter()
	token := signMapClaimsToken(t, testJWTSecret, jwt.MapClaims{
		"user_id":  float64(42),
		"username": "admin_one",
		"role":     "admin",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/probe", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}
	if w.Body.String() != "uid=42" {
		t.Fatalf("body = %q, want uid=42", w.Body.String())
	}
}

func TestAuthMiddlewareRejectsStudentWith403(t *testing.T) {
	t.Setenv("JWT_SECRET", testJWTSecret)
	r := newAuthProbeRouter()
	token := signMapClaimsToken(t, testJWTSecret, jwt.MapClaims{
		"user_id":  float64(7),
		"username": "student_one",
		"role":     "student",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/probe", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403, body = %s", w.Code, w.Body.String())
	}
}

func TestAuthMiddlewareAllowsPlatformAdmin(t *testing.T) {
	t.Setenv("JWT_SECRET", testJWTSecret)
	r := newAuthProbeRouter()
	token := signMapClaimsToken(t, testJWTSecret, jwt.MapClaims{
		"user_id":  float64(1),
		"username": "root_admin",
		"role":     "platform_admin",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/probe", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}
}

