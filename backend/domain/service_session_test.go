package domain

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestSweepExpiredSessions(t *testing.T) {
	synctest.Test(
		t, func(t *testing.T) {
			s := &service{sessions: make(sessionStore)}

			s.sessions["expired"] = Session{
				SessionID: "expired",
				UserID:    1,
				ExpiresAt: time.Now().Add(-time.Hour),
			}
			s.sessions["valid"] = Session{
				SessionID: "valid",
				UserID:    2,
				ExpiresAt: time.Now().Add(time.Hour),
			}

			ctx, cancel := context.WithCancel(context.Background())
			go s.sweepExpiredSessions(ctx)

			synctest.Wait()
			time.Sleep(time.Minute)
			synctest.Wait()

			if _, ok := s.sessions["expired"]; ok {
				t.Error("expected expired session to be swept, but it was still present")
			}
			if _, ok := s.sessions["valid"]; !ok {
				t.Error("expected valid session to be kept, but it was swept")
			}

			cancel()
			synctest.Wait()
		},
	)
}

func TestSweepStopsOnContextCancel(t *testing.T) {
	synctest.Test(
		t, func(t *testing.T) {
			s := &service{sessions: make(sessionStore)}

			ctx, cancel := context.WithCancel(context.Background())
			go s.sweepExpiredSessions(ctx)

			synctest.Wait()
			cancel()
			synctest.Wait()
		},
	)
}
