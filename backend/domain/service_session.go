package domain

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var _ SessionService = &service{}

const (
	SessionDuration = 60 * 60 * 8 // 8 hours
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type Session struct {
	SessionID  string
	UserID     int64
	ExpiresAt  time.Time
	PrivateKey *rsa.PrivateKey
}

type sessionStore map[string]Session

func (s *service) CreateSession(ctx context.Context, userID int64) (Session, error) {
	u := generateUUID()
	session := Session{
		SessionID: u,
		UserID:    userID,
		ExpiresAt: time.Now().Add(SessionDuration * time.Second),
	}

	hasKeys, err := s.database.HasKeys(ctx, userID)
	if err != nil {
		return session, err
	}
	if hasKeys {
		keyPair, err := s.database.SelectKeys(ctx, userID)
		if err != nil {
			return session, err
		}
		if keyPair.PasswordID == nil {
			priv, err := x509.ParsePKCS1PrivateKey(keyPair.Private)
			if err != nil {
				return session, err
			}
			session.PrivateKey = priv
		}
	}

	s.sessions[u] = session
	return session, nil
}

func (s *service) GetSession(token string) (*Session, error) {
	t := strings.ToLower(token)
	if session, ok := s.sessions[t]; ok {
		if time.Now().After(session.ExpiresAt) {
			_ = s.DeleteSession(token)
			return nil, ErrSessionNotFound
		}
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
