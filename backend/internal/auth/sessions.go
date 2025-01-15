package auth

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

const (
	SessionDuration = 60 * 60 * 8 // 8 hours
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type SessionService interface {
	// CreateSession creates a new session for the given user.
	CreateSession(user int32) (*Session, error)

	// GetSession returns the session for the given token.
	GetSession(token string) (*Session, error)

	// DeleteSession deletes the session for the given token.
	DeleteSession(token string) error
}

type Session struct {
	SessionID string
	UserID    int32
	ExpiresAt time.Time
}

type sessionStore map[string]*Session

type sessionService struct {
	store sessionStore
}

func NewSessionService() SessionService {
	return &sessionService{
		store: make(sessionStore),
	}
}

func (s *sessionService) CreateSession(user int32) (*Session, error) {
	u := generateUUID()
	session := &Session{
		SessionID: u,
		UserID:    user,
		ExpiresAt: time.Now().Add(SessionDuration * time.Second),
	}
	s.store[u] = session
	return session, nil
}

func (s *sessionService) GetSession(token string) (*Session, error) {
	if session, ok := s.store[token]; ok {
		return session, nil
	}
	return nil, ErrSessionNotFound
}

func (s *sessionService) DeleteSession(token string) error {
	delete(s.store, token)
	return nil
}

func generateUUID() string {
	return uuid.NewString()
}
