package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/torfstack/synod/backend/domain"
)

func TestUnsealRateLimiter(t *testing.T) {
	e := echo.New()

	injectSession := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			setSession(c, &domain.Session{UserID: 42})
			return next(c)
		}
	}

	e.POST(
		"/unseal", func(c echo.Context) error {
			return c.NoContent(http.StatusNoContent)
		}, injectSession, newUnsealRateLimiter(),
	)

	// The first 5 requests should pass (burst of 5).
	for i := range 5 {
		req := httptest.NewRequest(http.MethodPost, "/unseal", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("request %d: expected %d, got %d", i+1, http.StatusNoContent, rec.Code)
		}
	}

	// 6th request should be rate-limited.
	req := httptest.NewRequest(http.MethodPost, "/unseal", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("request 6: expected %d, got %d", http.StatusTooManyRequests, rec.Code)
	}
}

func TestUnsealRateLimiterIsolatedPerUser(t *testing.T) {
	e := echo.New()

	makeInjectSession := func(userID int64) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				setSession(c, &domain.Session{UserID: userID})
				return next(c)
			}
		}
	}

	rateLimiter := newUnsealRateLimiter()
	e.POST(
		"/unseal/user1", func(c echo.Context) error {
			return c.NoContent(http.StatusNoContent)
		}, makeInjectSession(1), rateLimiter,
	)
	e.POST(
		"/unseal/user2", func(c echo.Context) error {
			return c.NoContent(http.StatusNoContent)
		}, makeInjectSession(2), rateLimiter,
	)

	// Exhaust the limit for user 1.
	for range 5 {
		req := httptest.NewRequest(http.MethodPost, "/unseal/user1", nil)
		e.ServeHTTP(httptest.NewRecorder(), req)
	}

	// User 1 should now be rate limited.
	req := httptest.NewRequest(http.MethodPost, "/unseal/user1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("user1 request 6: expected %d, got %d", http.StatusTooManyRequests, rec.Code)
	}

	// User 2 should still be allowed (separate bucket).
	req = httptest.NewRequest(http.MethodPost, "/unseal/user2", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Errorf("user2 request 1: expected %d, got %d", http.StatusNoContent, rec.Code)
	}
}
