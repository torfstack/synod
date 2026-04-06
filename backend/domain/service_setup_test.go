package domain

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/models"
)

// ---- IsUserSetup ----

func TestIsUserSetup_DelegatesToDatabase(t *testing.T) {
	db := &mockDatabase{
		hasKeysFn: func(_ context.Context, userID int64) (bool, error) {
			return userID == 42, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	setup, err := svc.IsUserSetup(context.Background(), Session{UserID: 42})
	require.NoError(t, err)
	assert.True(t, setup)

	setup, err = svc.IsUserSetup(context.Background(), Session{UserID: 99})
	require.NoError(t, err)
	assert.False(t, setup)
}

func TestIsUserSetup_PropagatesDatabaseError(t *testing.T) {
	dbErr := errors.New("db error")
	db := &mockDatabase{
		hasKeysFn: func(_ context.Context, _ int64) (bool, error) {
			return false, dbErr
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	_, err := svc.IsUserSetup(context.Background(), Session{UserID: 1})
	require.ErrorIs(t, err, dbErr)
}

// ---- SetupUserPlain ----

func TestSetupUserPlain_StoresCipherInSession(t *testing.T) {
	sessionID := "test-session"
	db := &mockDatabase{
		insertKeysFn: func(_ context.Context, kp models.UserKeyPair) (models.UserKeyPair, error) {
			return kp, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}
	session := Session{SessionID: sessionID, UserID: 1}

	err := svc.SetupUserPlain(context.Background(), session)

	require.NoError(t, err)
	stored, ok := svc.sessions[sessionID]
	require.True(t, ok, "session should be updated in store")
	assert.NotNil(t, stored.Cipher)
}

func TestSetupUserPlain_DatabaseInsertKeysFails_ReturnsError(t *testing.T) {
	dbErr := errors.New("insert failed")
	db := &mockDatabase{
		insertKeysFn: func(_ context.Context, _ models.UserKeyPair) (models.UserKeyPair, error) {
			return models.UserKeyPair{}, dbErr
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	err := svc.SetupUserPlain(context.Background(), Session{SessionID: "s1", UserID: 1})

	require.ErrorIs(t, err, dbErr)
}

func TestSetupUserPlain_KeyMaterialHasNoKDFPrefix(t *testing.T) {
	var insertedKeyMaterial []byte
	db := &mockDatabase{
		insertKeysFn: func(_ context.Context, kp models.UserKeyPair) (models.UserKeyPair, error) {
			insertedKeyMaterial = kp.KeyMaterial
			return kp, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	err := svc.SetupUserPlain(context.Background(), Session{SessionID: "s2", UserID: 2})
	require.NoError(t, err)

	// Plain setup must NOT start with the KDF salt prefix
	assert.False(
		t,
		slices.Equal(insertedKeyMaterial[:4], crypto.KDFSaltPrefix),
		"plain setup should not embed KDF salt prefix in key material",
	)
}

// ---- SetupUserWithPassword ----

func TestSetupUserWithPassword_StoresCipherInSession(t *testing.T) {
	sessionID := "pwd-session"
	db := &mockDatabase{
		insertPasswordFn: func(_ context.Context, pw models.HashedPassword) (models.HashedPassword, error) {
			id := int64(1)
			pw.ID = &id
			return pw, nil
		},
		insertKeysFn: func(_ context.Context, kp models.UserKeyPair) (models.UserKeyPair, error) {
			return kp, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}
	session := Session{SessionID: sessionID, UserID: 3}

	err := svc.SetupUserWithPassword(context.Background(), session, crypto.Password("correct-horse-battery"))

	require.NoError(t, err)
	stored, ok := svc.sessions[sessionID]
	require.True(t, ok, "session should be updated in store")
	assert.NotNil(t, stored.Cipher)
}

func TestSetupUserWithPassword_KeyMaterialHasKDFPrefix(t *testing.T) {
	var insertedKeyMaterial []byte
	db := &mockDatabase{
		insertPasswordFn: func(_ context.Context, pw models.HashedPassword) (models.HashedPassword, error) {
			id := int64(1)
			pw.ID = &id
			return pw, nil
		},
		insertKeysFn: func(_ context.Context, kp models.UserKeyPair) (models.UserKeyPair, error) {
			insertedKeyMaterial = kp.KeyMaterial
			return kp, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	err := svc.SetupUserWithPassword(context.Background(), Session{SessionID: "s3", UserID: 4}, crypto.Password("passw0rd"))
	require.NoError(t, err)

	require.GreaterOrEqual(t, len(insertedKeyMaterial), 4)
	assert.True(
		t,
		slices.Equal(insertedKeyMaterial[:4], crypto.KDFSaltPrefix),
		"password setup key material must start with KDF salt prefix",
	)
}

func TestSetupUserWithPassword_InsertPasswordFails_ReturnsError(t *testing.T) {
	dbErr := errors.New("password insert failed")
	db := &mockDatabase{
		insertPasswordFn: func(_ context.Context, _ models.HashedPassword) (models.HashedPassword, error) {
			return models.HashedPassword{}, dbErr
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	err := svc.SetupUserWithPassword(context.Background(), Session{SessionID: "s4", UserID: 5}, crypto.Password("pw"))

	require.ErrorIs(t, err, dbErr)
}

// ---- UnsealWithPassword ----

func TestUnsealWithPassword_AlreadyUnsealed_IsNoOp(t *testing.T) {
	cipher := newTestCipher(t)
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}
	session := &Session{SessionID: "sealed", UserID: 1, Cipher: cipher}

	err := svc.UnsealWithPassword(context.Background(), session, crypto.Password("any-password"))

	require.NoError(t, err)
	// cipher should be unchanged
	assert.Equal(t, cipher, session.Cipher)
}

func TestUnsealWithPassword_FullRoundTrip(t *testing.T) {
	// Simulate the full password setup + unseal cycle.
	const password = "correct-horse-battery-staple"
	sessionID := "round-trip"

	var storedKeys models.UserKeyPair
	var storedPassword models.HashedPassword
	var nextPasswordID int64 = 10

	db := &mockDatabase{
		insertPasswordFn: func(_ context.Context, pw models.HashedPassword) (models.HashedPassword, error) {
			storedPassword = pw
			storedPassword.ID = &nextPasswordID
			return storedPassword, nil
		},
		insertKeysFn: func(_ context.Context, kp models.UserKeyPair) (models.UserKeyPair, error) {
			storedKeys = kp
			return storedKeys, nil
		},
		selectKeysFn: func(_ context.Context, _ int64) (models.UserKeyPair, error) {
			return storedKeys, nil
		},
		selectPasswordFn: func(_ context.Context, _ int64) (models.HashedPassword, error) {
			return storedPassword, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	// Setup phase
	err := svc.SetupUserWithPassword(context.Background(), Session{SessionID: sessionID, UserID: 1}, crypto.Password(password))
	require.NoError(t, err)

	// Clear the cipher to simulate a sealed session (e.g. after server restart)
	svc.sessions[sessionID] = Session{SessionID: sessionID, UserID: 1, Cipher: nil}
	session := &Session{SessionID: sessionID, UserID: 1, Cipher: nil}

	// Unseal phase
	err = svc.UnsealWithPassword(context.Background(), session, crypto.Password(password))
	require.NoError(t, err)
	assert.NotNil(t, session.Cipher, "cipher should be set after successful unseal")
}

func TestUnsealWithPassword_WrongPassword_ReturnsError(t *testing.T) {
	const password = "correct-password"

	var storedKeys models.UserKeyPair
	var storedPassword models.HashedPassword
	var nextPasswordID int64 = 20

	db := &mockDatabase{
		insertPasswordFn: func(_ context.Context, pw models.HashedPassword) (models.HashedPassword, error) {
			storedPassword = pw
			storedPassword.ID = &nextPasswordID
			return storedPassword, nil
		},
		insertKeysFn: func(_ context.Context, kp models.UserKeyPair) (models.UserKeyPair, error) {
			storedKeys = kp
			return storedKeys, nil
		},
		selectKeysFn: func(_ context.Context, _ int64) (models.UserKeyPair, error) {
			return storedKeys, nil
		},
		selectPasswordFn: func(_ context.Context, _ int64) (models.HashedPassword, error) {
			return storedPassword, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	err := svc.SetupUserWithPassword(context.Background(), Session{SessionID: "s5", UserID: 2}, crypto.Password(password))
	require.NoError(t, err)

	session := &Session{SessionID: "s5", UserID: 2, Cipher: nil}
	err = svc.UnsealWithPassword(context.Background(), session, crypto.Password("wrong-password"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "hash mismatch")
}

func TestUnsealWithPassword_NoPasswordID_ReturnsError(t *testing.T) {
	db := &mockDatabase{
		selectKeysFn: func(_ context.Context, _ int64) (models.UserKeyPair, error) {
			// KeyPair without a PasswordID means it is not password-protected
			return models.UserKeyPair{PasswordID: nil}, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}
	session := &Session{SessionID: "s6", UserID: 3, Cipher: nil}

	err := svc.UnsealWithPassword(context.Background(), session, crypto.Password("pw"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no password")
}

func TestUnsealWithPassword_SelectKeysFails_ReturnsError(t *testing.T) {
	dbErr := errors.New("keys not found")
	db := &mockDatabase{
		selectKeysFn: func(_ context.Context, _ int64) (models.UserKeyPair, error) {
			return models.UserKeyPair{}, dbErr
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}
	session := &Session{SessionID: "s7", UserID: 4, Cipher: nil}

	err := svc.UnsealWithPassword(context.Background(), session, crypto.Password("pw"))

	require.ErrorIs(t, err, dbErr)
}

func TestUnsealWithPassword_UnsupportedKeyMaterial_ReturnsError(t *testing.T) {
	const password = "my-password"
	var storedPassword models.HashedPassword
	var nextPasswordID int64 = 30

	db := &mockDatabase{
		insertPasswordFn: func(_ context.Context, pw models.HashedPassword) (models.HashedPassword, error) {
			storedPassword = pw
			storedPassword.ID = &nextPasswordID
			return storedPassword, nil
		},
		selectKeysFn: func(_ context.Context, _ int64) (models.UserKeyPair, error) {
			// Deliberately provide key material without the expected KDF prefix
			return models.UserKeyPair{
				PasswordID:  &nextPasswordID,
				KeyMaterial: []byte("short"), // too short and wrong prefix
			}, nil
		},
		selectPasswordFn: func(_ context.Context, _ int64) (models.HashedPassword, error) {
			return storedPassword, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	// Hash and store the password first so the hash comparison passes
	hashedPw, err := crypto.HashPassword(crypto.Password(password))
	require.NoError(t, err)
	storedPassword = models.HashedPassword{
		Hash:       hashedPw.Hash,
		Salt:       hashedPw.Salt,
		Iterations: hashedPw.IterationsUsed,
		ID:         &nextPasswordID,
	}

	session := &Session{SessionID: "s8", UserID: 5, Cipher: nil}
	err = svc.UnsealWithPassword(context.Background(), session, crypto.Password(password))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported format")
}
