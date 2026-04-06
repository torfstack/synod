package domain

import (
	"context"

	"github.com/torfstack/synod/backend/db"
	"github.com/torfstack/synod/backend/models"
)

// mockDatabase is a test double for db.Database.
// Each field is a function so individual tests can override only what they need.
type mockDatabase struct {
	doesUserExistFn    func(ctx context.Context, username string) (bool, error)
	insertUserFn       func(ctx context.Context, user models.User) (models.ExistingUser, error)
	selectUserByNameFn func(ctx context.Context, username string) (models.ExistingUser, error)

	upsertSecretFn  func(ctx context.Context, secret models.EncryptedSecret, userID int64) (models.EncryptedSecret, error)
	selectSecretsFn func(ctx context.Context, userID int64) ([]models.EncryptedSecret, error)

	insertKeysFn func(ctx context.Context, pair models.UserKeyPair) (models.UserKeyPair, error)
	selectKeysFn func(ctx context.Context, userID int64) (models.UserKeyPair, error)
	hasKeysFn    func(ctx context.Context, userID int64) (bool, error)

	insertPasswordFn func(ctx context.Context, password models.HashedPassword) (models.HashedPassword, error)
	selectPasswordFn func(ctx context.Context, passwordID int64) (models.HashedPassword, error)

	withTxFn func(ctx context.Context, fn func(db.Database) error) error
}

var _ db.Database = (*mockDatabase)(nil)

func (m *mockDatabase) WithTx(ctx context.Context, fn func(db.Database) error) error {
	if m.withTxFn != nil {
		return m.withTxFn(ctx, fn)
	}
	// Default: run the function directly using self
	return fn(m)
}

func (m *mockDatabase) DoesUserExist(ctx context.Context, username string) (bool, error) {
	if m.doesUserExistFn != nil {
		return m.doesUserExistFn(ctx, username)
	}
	return false, nil
}

func (m *mockDatabase) InsertUser(ctx context.Context, user models.User) (models.ExistingUser, error) {
	if m.insertUserFn != nil {
		return m.insertUserFn(ctx, user)
	}
	return models.ExistingUser{}, nil
}

func (m *mockDatabase) SelectUserByName(ctx context.Context, username string) (models.ExistingUser, error) {
	if m.selectUserByNameFn != nil {
		return m.selectUserByNameFn(ctx, username)
	}
	return models.ExistingUser{}, nil
}

func (m *mockDatabase) UpsertSecret(ctx context.Context, secret models.EncryptedSecret, userID int64) (models.EncryptedSecret, error) {
	if m.upsertSecretFn != nil {
		return m.upsertSecretFn(ctx, secret, userID)
	}
	return secret, nil
}

func (m *mockDatabase) SelectSecrets(ctx context.Context, userID int64) ([]models.EncryptedSecret, error) {
	if m.selectSecretsFn != nil {
		return m.selectSecretsFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockDatabase) InsertKeys(ctx context.Context, pair models.UserKeyPair) (models.UserKeyPair, error) {
	if m.insertKeysFn != nil {
		return m.insertKeysFn(ctx, pair)
	}
	return pair, nil
}

func (m *mockDatabase) SelectKeys(ctx context.Context, userID int64) (models.UserKeyPair, error) {
	if m.selectKeysFn != nil {
		return m.selectKeysFn(ctx, userID)
	}
	return models.UserKeyPair{}, nil
}

func (m *mockDatabase) HasKeys(ctx context.Context, userID int64) (bool, error) {
	if m.hasKeysFn != nil {
		return m.hasKeysFn(ctx, userID)
	}
	return false, nil
}

func (m *mockDatabase) InsertPassword(ctx context.Context, password models.HashedPassword) (models.HashedPassword, error) {
	if m.insertPasswordFn != nil {
		return m.insertPasswordFn(ctx, password)
	}
	return password, nil
}

func (m *mockDatabase) SelectPassword(ctx context.Context, passwordID int64) (models.HashedPassword, error) {
	if m.selectPasswordFn != nil {
		return m.selectPasswordFn(ctx, passwordID)
	}
	return models.HashedPassword{}, nil
}
