package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/torfstack/synod/backend/domain"
)

// ---- IsAuthorized ----

func TestIsAuthorized_NoSessionCookie_Returns401(t *testing.T) {
	svc := &mockDomainService{}
	s := NewServer(testConfig(), svc)
	e := echo.New()
	e.GET("/auth", s.IsAuthorized)

	req := httptest.NewRequest(http.MethodGet, "/auth", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestIsAuthorized_SessionNotFound_Returns401WithJSON(t *testing.T) {
	svc := &mockDomainService{
		getSessionFn: func(token string) (*domain.Session, error) {
			return nil, domain.ErrSessionNotFound
		},
	}
	s := NewServer(testConfig(), svc)
	e := echo.New()
	e.GET("/auth", s.IsAuthorized)

	req := httptest.NewRequest(http.MethodGet, "/auth", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "bad-token"})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "isAuthenticated")
}

func TestIsAuthorized_ValidSession_Returns200WithAuthStatus(t *testing.T) {
	session := &domain.Session{
		SessionID: "valid-token",
		UserID:    42,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	svc := &mockDomainService{
		getSessionFn: func(token string) (*domain.Session, error) {
			if token == "valid-token" {
				return session, nil
			}
			return nil, domain.ErrSessionNotFound
		},
		isUserSetupFn: func(_ context.Context, _ domain.Session) (bool, error) {
			return true, nil
		},
	}
	s := NewServer(testConfig(), svc)
	e := echo.New()
	e.GET("/auth", s.IsAuthorized)

	req := httptest.NewRequest(http.MethodGet, "/auth", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "valid-token"})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `"isAuthenticated":true`)
	assert.Contains(t, body, `"isSetup":true`)
}

func TestIsAuthorized_SealedSession_NeedsToUnsealIsTrue(t *testing.T) {
	// Cipher is nil -> session is sealed
	session := &domain.Session{
		SessionID: "sealed-token",
		UserID:    1,
		ExpiresAt: time.Now().Add(time.Hour),
		Cipher:    nil,
	}
	svc := &mockDomainService{
		getSessionFn: func(_ string) (*domain.Session, error) {
			return session, nil
		},
		isUserSetupFn: func(_ context.Context, _ domain.Session) (bool, error) {
			return true, nil
		},
	}
	s := NewServer(testConfig(), svc)
	e := echo.New()
	e.GET("/auth", s.IsAuthorized)

	req := httptest.NewRequest(http.MethodGet, "/auth", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "sealed-token"})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"needsToUnseal":true`)
}

// ---- EndSession ----

func TestEndSession_NoSessionCookie_Returns200(t *testing.T) {
	svc := &mockDomainService{}
	s := NewServer(testConfig(), svc)
	e := echo.New()
	e.DELETE("/auth", s.EndSession)

	req := httptest.NewRequest(http.MethodDelete, "/auth", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestEndSession_WithSessionCookie_DeletesSessionAndClears(t *testing.T) {
	deletedToken := ""
	svc := &mockDomainService{
		deleteSessionFn: func(token string) error {
			deletedToken = token
			return nil
		},
	}
	s := NewServer(testConfig(), svc)
	e := echo.New()
	e.DELETE("/auth", s.EndSession)

	req := httptest.NewRequest(http.MethodDelete, "/auth", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "my-session-id"})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "my-session-id", deletedToken)

	// Confirm an empty/expired session cookie is set
	found := false
	for _, c := range rec.Result().Cookies() {
		if c.Name == SessionCookieName {
			assert.Equal(t, "", c.Value)
			found = true
		}
	}
	assert.True(t, found, "response should clear the session cookie")
}
