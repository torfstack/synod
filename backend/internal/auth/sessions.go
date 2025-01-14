package auth

import (
    "errors"
    "time"
    "github.com/google/uuid"
)

const (
    SessionDuration = 60 * 60 // 1 hour
)

var (
    ErrSessionNotFound = errors.New("session not found")
)

type SessionService interface {
    // CreateSession creates a new session for the given user.
    CreateSession(user int64) (*Session, error)

    // GetSession returns the session for the given token.
    GetSession(token string) (*Session, error)

    // DeleteSession deletes the session for the given token.
    DeleteSession(token string) error
}

type Session struct {
    Token string
    UserID int64
    ExpiresAt int64
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

func (s *sessionService) CreateSession(user int64) (*Session, error) {
    uuid := generateUUID()
    session := &Session{
        Token: uuid,
        UserID: user,
        ExpiresAt: time.Now().Unix() + SessionDuration,
    }
    s.store[uuid] = session
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
