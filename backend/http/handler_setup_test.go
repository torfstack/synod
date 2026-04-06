package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/domain"
)

// ---- UnsealWithPassword ----

func TestUnsealWithPassword_NoSession_Returns401(t *testing.T) {
	e := echo.New()
	s := NewServer(testConfig(), &mockDomainService{})
	e.POST("/unseal", s.UnsealWithPassword)

	req := httptest.NewRequest(http.MethodPost, "/unseal", strings.NewReader(`{"password":"pw"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUnsealWithPassword_InvalidBody_Returns400(t *testing.T) {
	session := &domain.Session{UserID: 1}
	s := NewServer(testConfig(), &mockDomainService{})
	e := newEchoWithSession(session)
	e.POST("/unseal", s.UnsealWithPassword)

	req := httptest.NewRequest(http.MethodPost, "/unseal", strings.NewReader(`not-json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUnsealWithPassword_WrongPassword_Returns403(t *testing.T) {
	session := &domain.Session{UserID: 1}
	svc := &mockDomainService{
		unsealWithPasswordFn: func(_ context.Context, _ *domain.Session, _ crypto.Password) error {
			return domain.ErrInvalidPassword
		},
	}
	s := NewServer(testConfig(), svc)
	e := newEchoWithSession(session)
	e.POST("/unseal", s.UnsealWithPassword)

	req := httptest.NewRequest(http.MethodPost, "/unseal", strings.NewReader(`{"password":"wrong"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestUnsealWithPassword_Success_Returns204(t *testing.T) {
	session := &domain.Session{UserID: 1}
	svc := &mockDomainService{
		unsealWithPasswordFn: func(_ context.Context, _ *domain.Session, _ crypto.Password) error {
			return nil
		},
	}
	s := NewServer(testConfig(), svc)
	e := newEchoWithSession(session)
	e.POST("/unseal", s.UnsealWithPassword)

	req := httptest.NewRequest(http.MethodPost, "/unseal", strings.NewReader(`{"password":"correct"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// ---- PostSetupPlain ----

func TestPostSetupPlain_NoSession_Returns401(t *testing.T) {
	e := echo.New()
	s := NewServer(testConfig(), &mockDomainService{})
	e.POST("/setup/plain", s.PostSetupPlain)

	req := httptest.NewRequest(http.MethodPost, "/setup/plain", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestPostSetupPlain_ServiceError_PropagatesError(t *testing.T) {
	session := &domain.Session{UserID: 1}
	svc := &mockDomainService{
		setupUserPlainFn: func(_ context.Context, _ domain.Session) error {
			return errors.New("setup failed")
		},
	}
	s := NewServer(testConfig(), svc)
	e := newEchoWithSession(session)
	e.POST("/setup/plain", s.PostSetupPlain)

	req := httptest.NewRequest(http.MethodPost, "/setup/plain", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.NotEqual(t, http.StatusCreated, rec.Code)
}

func TestPostSetupPlain_Success_Returns201(t *testing.T) {
	session := &domain.Session{UserID: 1}
	svc := &mockDomainService{
		setupUserPlainFn: func(_ context.Context, _ domain.Session) error {
			return nil
		},
	}
	s := NewServer(testConfig(), svc)
	e := newEchoWithSession(session)
	e.POST("/setup/plain", s.PostSetupPlain)

	req := httptest.NewRequest(http.MethodPost, "/setup/plain", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

// ---- PostSetupPassword ----

func TestPostSetupPassword_NoSession_Returns401(t *testing.T) {
	e := echo.New()
	s := NewServer(testConfig(), &mockDomainService{})
	e.POST("/setup/password", s.PostSetupPassword)

	req := httptest.NewRequest(http.MethodPost, "/setup/password", strings.NewReader(`{"password":"pw"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestPostSetupPassword_InvalidBody_Returns400(t *testing.T) {
	session := &domain.Session{UserID: 1}
	s := NewServer(testConfig(), &mockDomainService{})
	e := newEchoWithSession(session)
	e.POST("/setup/password", s.PostSetupPassword)

	req := httptest.NewRequest(http.MethodPost, "/setup/password", strings.NewReader(`not-json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPostSetupPassword_ServiceError_PropagatesError(t *testing.T) {
	session := &domain.Session{UserID: 1}
	svc := &mockDomainService{
		setupUserWithPasswordFn: func(_ context.Context, _ domain.Session, _ crypto.Password) error {
			return errors.New("setup failed")
		},
	}
	s := NewServer(testConfig(), svc)
	e := newEchoWithSession(session)
	e.POST("/setup/password", s.PostSetupPassword)

	req := httptest.NewRequest(http.MethodPost, "/setup/password", strings.NewReader(`{"password":"secure-pass"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.NotEqual(t, http.StatusCreated, rec.Code)
}

func TestPostSetupPassword_Success_Returns201(t *testing.T) {
	session := &domain.Session{UserID: 1}
	svc := &mockDomainService{
		setupUserWithPasswordFn: func(_ context.Context, _ domain.Session, _ crypto.Password) error {
			return nil
		},
	}
	s := NewServer(testConfig(), svc)
	e := newEchoWithSession(session)
	e.POST("/setup/password", s.PostSetupPassword)

	req := httptest.NewRequest(http.MethodPost, "/setup/password", strings.NewReader(`{"password":"secure-pass"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}
