package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/torfstack/synod/backend/models"
)

// ---- CreateSession ----

func TestCreateSession_NoKeys_SessionHasNoCipher(t *testing.T) {
	db := &mockDatabase{
		hasKeysFn: func(_ context.Context, _ int64) (bool, error) {
			return false, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	session, err := svc.CreateSession(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), session.UserID)
	assert.NotEmpty(t, session.SessionID)
	assert.Nil(t, session.Cipher)
	assert.True(t, session.ExpiresAt.After(time.Now()))

	// Session should be stored
	stored, ok := svc.sessions[session.SessionID]
	require.True(t, ok)
	assert.Equal(t, session.UserID, stored.UserID)
}

func TestCreateSession_HasPlainKeys_SessionHasCipher(t *testing.T) {
	cipher := newTestCipher(t)
	keyBytes, err := cipher.Serialize()
	require.NoError(t, err)

	db := &mockDatabase{
		hasKeysFn: func(_ context.Context, _ int64) (bool, error) {
			return true, nil
		},
		selectKeysFn: func(_ context.Context, _ int64) (models.UserKeyPair, error) {
			return models.UserKeyPair{KeyMaterial: keyBytes, PasswordID: nil}, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	session, err := svc.CreateSession(context.Background(), 2)

	require.NoError(t, err)
	assert.NotNil(t, session.Cipher)
}

func TestCreateSession_HasPasswordProtectedKeys_SessionHasNoCipher(t *testing.T) {
	passwordID := int64(1)
	db := &mockDatabase{
		hasKeysFn: func(_ context.Context, _ int64) (bool, error) {
			return true, nil
		},
		selectKeysFn: func(_ context.Context, _ int64) (models.UserKeyPair, error) {
			return models.UserKeyPair{PasswordID: &passwordID}, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	session, err := svc.CreateSession(context.Background(), 3)

	require.NoError(t, err)
	assert.Nil(t, session.Cipher, "password-protected keys should not load the cipher at login")
}

func TestCreateSession_HasKeysDatabaseError_ReturnsError(t *testing.T) {
	dbErr := errors.New("select keys failed")
	db := &mockDatabase{
		hasKeysFn: func(_ context.Context, _ int64) (bool, error) {
			return true, nil
		},
		selectKeysFn: func(_ context.Context, _ int64) (models.UserKeyPair, error) {
			return models.UserKeyPair{}, dbErr
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	_, err := svc.CreateSession(context.Background(), 4)

	require.ErrorIs(t, err, dbErr)
}

func TestCreateSession_HasKeysFails_ReturnsError(t *testing.T) {
	dbErr := errors.New("has keys failed")
	db := &mockDatabase{
		hasKeysFn: func(_ context.Context, _ int64) (bool, error) {
			return false, dbErr
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	_, err := svc.CreateSession(context.Background(), 5)
	require.ErrorIs(t, err, dbErr)
}

func TestCreateSession_GeneratesUniqueSessionIDs(t *testing.T) {
	db := &mockDatabase{
		hasKeysFn: func(_ context.Context, _ int64) (bool, error) { return false, nil },
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	s1, err1 := svc.CreateSession(context.Background(), 1)
	s2, err2 := svc.CreateSession(context.Background(), 1)

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, s1.SessionID, s2.SessionID)
}

// ---- GetSession ----

func TestGetSession_ValidSession_ReturnsSession(t *testing.T) {
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}
	svc.sessions["abc"] = Session{SessionID: "abc", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}

	session, err := svc.GetSession("abc")

	require.NoError(t, err)
	assert.Equal(t, int64(1), session.UserID)
}

func TestGetSession_UnknownToken_ReturnsErrSessionNotFound(t *testing.T) {
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}

	_, err := svc.GetSession("no-such-token")

	require.ErrorIs(t, err, ErrSessionNotFound)
}

func TestGetSession_ExpiredSession_ReturnsErrSessionNotFoundAndDeletesIt(t *testing.T) {
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}
	svc.sessions["expired"] = Session{
		SessionID: "expired",
		UserID:    1,
		ExpiresAt: time.Now().Add(-time.Minute),
	}

	_, err := svc.GetSession("expired")

	require.ErrorIs(t, err, ErrSessionNotFound)
	_, stillPresent := svc.sessions["expired"]
	assert.False(t, stillPresent, "expired session should be deleted on access")
}

func TestGetSession_TokenIsCaseInsensitive(t *testing.T) {
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}
	svc.sessions["abcdef"] = Session{SessionID: "abcdef", UserID: 7, ExpiresAt: time.Now().Add(time.Hour)}

	session, err := svc.GetSession("ABCDEF")

	require.NoError(t, err)
	assert.Equal(t, int64(7), session.UserID)
}

// ---- DeleteSession ----

func TestDeleteSession_RemovesSession(t *testing.T) {
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}
	svc.sessions["tok"] = Session{SessionID: "tok", UserID: 1}

	err := svc.DeleteSession("tok")

	require.NoError(t, err)
	_, ok := svc.sessions["tok"]
	assert.False(t, ok)
}

func TestDeleteSession_NonExistentToken_ReturnsNoError(t *testing.T) {
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}

	err := svc.DeleteSession("ghost")

	require.NoError(t, err)
}

func TestDeleteSession_TokenIsCaseInsensitive(t *testing.T) {
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}
	svc.sessions["token123"] = Session{SessionID: "token123", UserID: 1}

	err := svc.DeleteSession("TOKEN123")

	require.NoError(t, err)
	_, ok := svc.sessions["token123"]
	assert.False(t, ok)
}
