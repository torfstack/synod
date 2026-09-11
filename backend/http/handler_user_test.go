package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/torfstack/synod/backend/domain"
)

func TestLookUpUser_NoSearchString_Returns400(t *testing.T) {
	s := NewServer(testConfig(), &mockDomainService{})
	e := newEchoWithSession(&domain.Session{UserID: 1})
	e.GET("/users/lookup", s.LookUpUser)

	req := httptest.NewRequest(http.MethodGet, "/users/lookup", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "error")
}

func TestLookUpUser_WithSearchString_Returns200(t *testing.T) {
	s := NewServer(testConfig(), &mockDomainService{})
	e := newEchoWithSession(&domain.Session{UserID: 1})
	e.GET("/users/lookup", s.LookUpUser)

	req := httptest.NewRequest(http.MethodGet, "/users/lookup?find=alice", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
