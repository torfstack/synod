package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/domain"
	"github.com/torfstack/synod/backend/models"
)

// newEchoWithSession returns an Echo instance with a GET /test route that
// injects session into context before calling the handler.
func newEchoWithSession(session *domain.Session) *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if session != nil {
				setSession(c, session)
			}
			return next(c)
		}
	})
	return e
}

func newTestCipher(t *testing.T) *crypto.AsymmetricCipher {
	t.Helper()
	c, err := crypto.NewAsymmetricCipher()
	require.NoError(t, err)
	return c
}

// ---- GetSecrets ----

func TestGetSecrets_NoSession_Returns401(t *testing.T) {
	e := echo.New()
	svc := &mockDomainService{}
	s := NewServer(testConfig(), svc)

	e.GET("/secrets", s.GetSecrets)
	req := httptest.NewRequest(http.MethodGet, "/secrets", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetSecrets_ServiceError_PropagatesError(t *testing.T) {
	cipher := newTestCipher(t)
	session := &domain.Session{UserID: 1, Cipher: cipher}
	svc := &mockDomainService{
		getSecretsFn: func(_ context.Context, _ int64, _ *crypto.AsymmetricCipher) ([]models.Secret, error) {
			return nil, errors.New("db error")
		},
	}
	s := NewServer(testConfig(), svc)
	e := newEchoWithSession(session)
	e.GET("/secrets", s.GetSecrets)

	req := httptest.NewRequest(http.MethodGet, "/secrets", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.NotEqual(t, http.StatusOK, rec.Code)
}

func TestGetSecrets_Success_Returns200WithJSON(t *testing.T) {
	cipher := newTestCipher(t)
	session := &domain.Session{UserID: 1, Cipher: cipher}
	expected := []models.Secret{
		{Value: "v1", Key: "k1", Url: "u1", Tags: []string{"a"}},
	}
	svc := &mockDomainService{
		getSecretsFn: func(_ context.Context, _ int64, _ *crypto.AsymmetricCipher) ([]models.Secret, error) {
			return expected, nil
		},
	}
	s := NewServer(testConfig(), svc)
	e := newEchoWithSession(session)
	e.GET("/secrets", s.GetSecrets)

	req := httptest.NewRequest(http.MethodGet, "/secrets", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got []models.Secret
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Len(t, got, 1)
	assert.Equal(t, "v1", got[0].Value)
}

// ---- PostSecret ----

func TestPostSecret_NoSession_Returns401(t *testing.T) {
	e := echo.New()
	s := NewServer(testConfig(), &mockDomainService{})
	e.POST("/secrets", s.PostSecret)

	req := httptest.NewRequest(
		http.MethodPost,
		"/secrets",
		strings.NewReader(`{"value":"v","key":"k","url":"u","tags":[]}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestPostSecret_InvalidBody_Returns400(t *testing.T) {
	cipher := newTestCipher(t)
	session := &domain.Session{UserID: 1, Cipher: cipher}
	s := NewServer(testConfig(), &mockDomainService{})
	e := newEchoWithSession(session)
	e.POST("/secrets", s.PostSecret)

	req := httptest.NewRequest(http.MethodPost, "/secrets", strings.NewReader(`not-json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPostSecret_ServiceError_PropagatesError(t *testing.T) {
	cipher := newTestCipher(t)
	session := &domain.Session{UserID: 1, Cipher: cipher}
	svc := &mockDomainService{
		upsertSecretFn: func(_ context.Context, _ models.Secret, _ int64, _ *crypto.AsymmetricCipher) (models.EncryptedSecret, error) {
			return models.EncryptedSecret{}, errors.New("upsert failed")
		},
	}
	s := NewServer(testConfig(), svc)
	e := newEchoWithSession(session)
	e.POST("/secrets", s.PostSecret)

	req := httptest.NewRequest(
		http.MethodPost,
		"/secrets",
		strings.NewReader(`{"value":"v","key":"k","url":"u","tags":[]}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.NotEqual(t, http.StatusCreated, rec.Code)
}

func TestPostSecret_Success_Returns201(t *testing.T) {
	cipher := newTestCipher(t)
	session := &domain.Session{UserID: 1, Cipher: cipher}
	svc := &mockDomainService{
		upsertSecretFn: func(_ context.Context, _ models.Secret, _ int64, _ *crypto.AsymmetricCipher) (models.EncryptedSecret, error) {
			return models.EncryptedSecret{Value: "encrypted", Key: "k"}, nil
		},
	}
	s := NewServer(testConfig(), svc)
	e := newEchoWithSession(session)
	e.POST("/secrets", s.PostSecret)

	req := httptest.NewRequest(
		http.MethodPost,
		"/secrets",
		strings.NewReader(`{"value":"my-secret","key":"k","url":"u","tags":[]}`),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}
