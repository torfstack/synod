package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/torfstack/synod/backend/crypto"
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
	UserID    int64
	ExpiresAt time.Time
	Cipher    *crypto.AsymmetricCipher
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
		key, err := s.database.SelectKeys(ctx, userID)
		if err != nil {
			return session, err
		}
		if key.PasswordID == nil {
			session.Cipher, err = crypto.AsymmetricCipherFromBytes(key.KeyMaterial)
			if err != nil {
				return session, err
			}
			if len(key.PublicKey) == 0 {
				publicKey, err := session.Cipher.SerializePublicKey()
				if err != nil {
					return session, err
				}
				if err := s.database.UpdatePublicKey(ctx, userID, publicKey); err != nil {
					return session, err
				}
			}
		}
	}

	s.sessionsMu.Lock()
	s.sessions[u] = session
	s.sessionsMu.Unlock()
	return session, nil
}

func (s *service) GetSession(token string) (*Session, error) {
	t := strings.ToLower(token)
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	if session, ok := s.sessions[t]; ok {
		if time.Now().After(session.ExpiresAt) {
			delete(s.sessions, t)
			return nil, ErrSessionNotFound
		}
		return &session, nil
	}
	return nil, ErrSessionNotFound
}

func (s *service) DeleteSession(token string) error {
	t := strings.ToLower(token)
	s.sessionsMu.Lock()
	delete(s.sessions, t)
	s.sessionsMu.Unlock()
	return nil
}

func generateUUID() string {
	return strings.ToLower(uuid.NewString())
}
