package domain

import (
	"context"
	"time"

	"github.com/torfstack/synod/backend/db"
)

type service struct {
	database db.Database
	sessions sessionStore
}

var _ Service = (*service)(nil)

func NewDomainService(ctx context.Context, db db.Database) Service {
	s := &service{
		database: db,
		sessions: make(sessionStore),
	}
	go s.sweepExpiredSessions(ctx)
	return s
}

func (s *service) sweepExpiredSessions(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			for id, session := range s.sessions {
				if now.After(session.ExpiresAt) {
					delete(s.sessions, id)
				}
			}
		}
	}
}
