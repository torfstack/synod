package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLookUpUser_NoSearchString_Returns400(t *testing.T) {
	s := NewServer(testConfig(), &mockDomainService{})
	e := echo.New()
	e.GET("/users/lookup", s.LookUpUser)

	req := httptest.NewRequest(http.MethodGet, "/users/lookup", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "error")
}

func TestLookUpUser_WithSearchString_Returns200(t *testing.T) {
	s := NewServer(testConfig(), &mockDomainService{})
	e := echo.New()
	e.GET("/users/lookup", s.LookUpUser)

	req := httptest.NewRequest(http.MethodGet, "/users/lookup?find=alice", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
