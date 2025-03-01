package domain

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

var _ SessionService = &service{}

const (
	SessionDuration = 60 * 60 * 8 // 8 hours
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type Session struct {
	SessionID string
	UserID    int32
	ExpiresAt time.Time
}

type sessionStore map[string]Session

func (s *service) CreateSession(user int32) (*Session, error) {
	u := generateUUID()
	session := Session{
		SessionID: u,
		UserID:    user,
		ExpiresAt: time.Now().Add(SessionDuration * time.Second),
	}
	s.sessions[u] = session
	return &session, nil
}

func (s *service) GetSession(token string) (*Session, error) {
	t := strings.ToLower(token)
	if session, ok := s.sessions[t]; ok {
		return &session, nil
	}
	return nil, ErrSessionNotFound
}

func (s *service) DeleteSession(token string) error {
	t := strings.ToLower(token)
	delete(s.sessions, t)
	return nil
}

func generateUUID() string {
	return strings.ToLower(uuid.NewString())
}
