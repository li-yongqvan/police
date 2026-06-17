package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"ai-forum/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	qqProvider                  = "qq"
	qqOAuthStateTTL             = 10 * time.Minute
	qqOAuthStateQueryParam      = "state"
	qqOAuthCodeQueryParam       = "code"
	qqOAuthReturnToQueryParam   = "return_to"
	qqOAuthDefaultReturnToPath  = "/community"
	qqOAuthAuthorizeEndpoint    = "https://graph.qq.com/oauth2.0/authorize"
	qqOAuthTokenEndpoint        = "https://graph.qq.com/oauth2.0/token"
	qqOAuthOpenIDEndpoint       = "https://graph.qq.com/oauth2.0/me"
	qqOAuthUserInfoEndpoint     = "https://graph.qq.com/user/get_user_info"
	envQQAppID                  = "QQ_APP_ID"
	envQQAppKey                 = "QQ_APP_KEY"
	envQQRedirectURI            = "QQ_REDIRECT_URI"
	envFrontendOAuthRedirectURL = "FRONTEND_OAUTH_REDIRECT_URL"
	envQQRequireInvite          = "QQ_OAUTH_REQUIRE_INVITE"
)

type QQOAuthHandler struct {
	Service *service.UserService
}

func NewQQOAuthHandler(svc *service.UserService) *QQOAuthHandler {
	return &QQOAuthHandler{Service: svc}
}

func (h *QQOAuthHandler) Start(c *gin.Context) {
	appID := strings.TrimSpace(os.Getenv(envQQAppID))
	redirectURI := strings.TrimSpace(os.Getenv(envQQRedirectURI))
	if appID == "" || redirectURI == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "QQ OAuth 未配置"})
		return
	}

	state, err := newRandomState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成登录状态失败"})
		return
	}

	returnTo := strings.TrimSpace(c.Query(qqOAuthReturnToQueryParam))
	if returnTo == "" {
		returnTo = qqOAuthDefaultReturnToPath
	}
	if err := h.Service.SetOAuthState(c.Request.Context(), qqProvider, state, returnTo, qqOAuthStateTTL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存登录状态失败"})
		return
	}

	u, _ := url.Parse(qqOAuthAuthorizeEndpoint)
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", appID)
	q.Set("redirect_uri", redirectURI)
	q.Set("state", state)
	q.Set("scope", "get_user_info")
	u.RawQuery = q.Encode()

	c.Redirect(http.StatusFound, u.String())
}

func (h *QQOAuthHandler) Callback(c *gin.Context) {
	code := strings.TrimSpace(c.Query(qqOAuthCodeQueryParam))
	state := strings.TrimSpace(c.Query(qqOAuthStateQueryParam))
	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 code 或 state"})
		return
	}

	returnTo, ok, err := h.Service.ConsumeOAuthState(c.Request.Context(), qqProvider, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "校验登录状态失败"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "登录状态已过期，请重试"})
		return
	}

	appID := strings.TrimSpace(os.Getenv(envQQAppID))
	appKey := strings.TrimSpace(os.Getenv(envQQAppKey))
	redirectURI := strings.TrimSpace(os.Getenv(envQQRedirectURI))
	if appID == "" || appKey == "" || redirectURI == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "QQ OAuth 未配置"})
		return
	}

	accessToken, err := qqExchangeCodeForToken(c.Request.Context(), appID, appKey, redirectURI, code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	openID, err := qqFetchOpenID(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	profile, _ := qqFetchUserInfo(c.Request.Context(), accessToken, appID, openID)

	requireInvite := strings.TrimSpace(os.Getenv(envQQRequireInvite)) == "1"
	u, role, tokens, err := h.Service.LoginOrCreateOAuthUser(c.Request.Context(), qqProvider, openID, profile, requireInvite)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	redirectBase := strings.TrimSpace(os.Getenv(envFrontendOAuthRedirectURL))
	if redirectBase == "" {
		redirectBase = "http://127.0.0.1:8091/oauth/qq"
	}
	target, err := buildFrontendRedirect(redirectBase, returnTo, tokens.AccessToken, tokens.RefreshToken, toUserJSON(u.ToResponse(), role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成跳转链接失败"})
		return
	}
	c.Redirect(http.StatusFound, target)
}

func newRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func qqExchangeCodeForToken(ctx context.Context, appID, appKey, redirectURI, code string) (string, error) {
	u, _ := url.Parse(qqOAuthTokenEndpoint)
	q := u.Query()
	q.Set("grant_type", "authorization_code")
	q.Set("client_id", appID)
	q.Set("client_secret", appKey)
	q.Set("code", code)
	q.Set("redirect_uri", redirectURI)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("QQ 登录失败，请稍后重试")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	s := strings.TrimSpace(string(body))
	if strings.Contains(s, "error") {
		return "", fmt.Errorf("QQ 登录失败，请重新授权")
	}

	vals, err := url.ParseQuery(s)
	if err != nil {
		return "", fmt.Errorf("QQ 登录失败，响应解析错误")
	}
	token := vals.Get("access_token")
	if token == "" {
		return "", fmt.Errorf("QQ 登录失败，未获取 access_token")
	}
	return token, nil
}

func qqFetchOpenID(ctx context.Context, accessToken string) (string, error) {
	u, _ := url.Parse(qqOAuthOpenIDEndpoint)
	q := u.Query()
	q.Set("access_token", accessToken)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("QQ 登录失败，请稍后重试")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	s := strings.TrimSpace(string(body))
	if strings.Contains(s, "error") {
		return "", fmt.Errorf("QQ 登录失败，未获取 openid")
	}

	// QQ returns: callback( {"client_id":"APPID","openid":"OPENID"} );
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || end <= start {
		return "", fmt.Errorf("QQ 登录失败，openid 响应格式错误")
	}
	var payload struct {
		OpenID string `json:"openid"`
	}
	if err := json.Unmarshal([]byte(s[start:end+1]), &payload); err != nil {
		return "", fmt.Errorf("QQ 登录失败，openid 解析错误")
	}
	if payload.OpenID == "" {
		return "", fmt.Errorf("QQ 登录失败，openid 为空")
	}
	return payload.OpenID, nil
}

func qqFetchUserInfo(ctx context.Context, accessToken, appID, openID string) (map[string]any, error) {
	u, _ := url.Parse(qqOAuthUserInfoEndpoint)
	q := u.Query()
	q.Set("access_token", accessToken)
	q.Set("oauth_consumer_key", appID)
	q.Set("openid", openID)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func buildFrontendRedirect(baseURL, returnTo, accessToken, refreshToken string, user map[string]any) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("token", accessToken)
	if refreshToken != "" {
		q.Set("refresh_token", refreshToken)
	}
	if returnTo != "" {
		q.Set("return_to", returnTo)
	}
	b, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	q.Set("user", base64.RawURLEncoding.EncodeToString(b))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

