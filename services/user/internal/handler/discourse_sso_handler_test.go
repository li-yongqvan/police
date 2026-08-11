package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleSSORedirectsAnonymousUserToPlatformLogin(t *testing.T) {
	t.Setenv("DISCOURSE_CONNECT_SECRET", "test-secret")
	gin.SetMode(gin.TestMode)

	payload := url.Values{
		"nonce":         {"nonce-1"},
		"return_sso_url": {"http://forum.test/session/sso_login"},
	}.Encode()
	rawPayload := base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(payload)))

	mac := hmac.New(sha256.New, []byte("test-secret"))
	_, _ = mac.Write([]byte(rawPayload))
	signature := hex.EncodeToString(mac.Sum(nil))

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/auth/discourse/sso?sso="+url.QueryEscape(rawPayload)+"&sig="+signature,
		nil,
	)
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request = request

	(&DiscourseSSOHandler{}).HandleSSO(context)

	if response.Code != http.StatusFound {
		t.Fatalf("expected status %d, got %d", http.StatusFound, response.Code)
	}
	if got := response.Header().Get("Location"); got != "/" {
		t.Fatalf("expected anonymous SSO redirect to /, got %q", got)
	}
}
