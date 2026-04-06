package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/torfstack/synod/backend/models"
)

// ---- DoesUserExist ----

func TestDoesUserExist_DelegatesToDatabase(t *testing.T) {
	db := &mockDatabase{
		doesUserExistFn: func(_ context.Context, username string) (bool, error) {
			return username == "alice", nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	exists, err := svc.DoesUserExist(context.Background(), "alice")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = svc.DoesUserExist(context.Background(), "unknown")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestDoesUserExist_PropagatesDatabaseError(t *testing.T) {
	dbErr := errors.New("connection refused")
	db := &mockDatabase{
		doesUserExistFn: func(_ context.Context, _ string) (bool, error) {
			return false, dbErr
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	_, err := svc.DoesUserExist(context.Background(), "alice")
	require.ErrorIs(t, err, dbErr)
}

// ---- InsertUser ----

func TestInsertUser_ReturnsCreatedUser(t *testing.T) {
	user := models.User{Subject: "sub-1", Email: "a@b.com", FullName: "Alice"}
	expected := models.ExistingUser{ID: 1, User: user}

	db := &mockDatabase{
		insertUserFn: func(_ context.Context, u models.User) (models.ExistingUser, error) {
			assert.Equal(t, user, u)
			return expected, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	got, err := svc.InsertUser(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestInsertUser_PropagatesDatabaseError(t *testing.T) {
	dbErr := errors.New("duplicate key")
	db := &mockDatabase{
		insertUserFn: func(_ context.Context, _ models.User) (models.ExistingUser, error) {
			return models.ExistingUser{}, dbErr
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	_, err := svc.InsertUser(context.Background(), models.User{})
	require.ErrorIs(t, err, dbErr)
}
