package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeRoleAuthority struct {
	role string
	err  error
}

func (f *fakeRoleAuthority) ResolveAuthoritativeRole(_ context.Context, _ uint) (string, error) {
	return f.role, f.err
}

func roleRouter(h *RoleHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/users/:id/role", h.GetUserRole)
	return r
}

func TestGetUserRoleReturnsAuthoritativeRole(t *testing.T) {
	h := &RoleHandler{authority: &fakeRoleAuthority{role: "admin"}}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/42/role", nil)
	roleRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body.UserID != 42 || body.Role != "admin" {
		t.Fatalf("body = %+v, want user_id=42 role=admin", body)
	}
}

func TestGetUserRoleUnknownUserResolvesToStudentWith200(t *testing.T) {
	// The Role Authority contract: no-assignment is not an error. The
	// handler must not turn the student fallback into a 404.
	h := &RoleHandler{authority: &fakeRoleAuthority{role: "student"}}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/7/role", nil)
	roleRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body struct {
		Role string `json:"role"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body.Role != "student" {
		t.Fatalf("body = %s (err %v), want role=student", w.Body.String(), err)
	}
}

func TestGetUserRoleRejectsInvalidID(t *testing.T) {
	h := &RoleHandler{authority: &fakeRoleAuthority{role: "admin"}}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/not-a-number/role", nil)
	roleRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestGetUserRoleAuthorityFailureReturns500(t *testing.T) {
	h := &RoleHandler{authority: &fakeRoleAuthority{err: errors.New("db down")}}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/42/role", nil)
	roleRouter(h).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}