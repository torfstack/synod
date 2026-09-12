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

	upsertSecretFn            func(ctx context.Context, secret models.EncryptedSecret, userID int64) (models.EncryptedSecret, error)
	selectSecretsFn           func(ctx context.Context, userID int64) ([]models.EncryptedSecret, error)
	selectAccessibleSecretsFn func(ctx context.Context, userID int64) ([]models.AccessibleSecret, error)
	selectSecretForOwnerFn    func(context.Context, int64, int64) (models.AccessibleSecret, error)
	insertSecretAccessFn      func(context.Context, int64, int64, int64, []byte) error
	selectUserBySharingIDFn   func(context.Context, string) (models.ExistingUser, error)
	selectPublicKeyFn         func(ctx context.Context, userID int64) ([]byte, error)
	insertThresholdSecretFn   func(context.Context, models.EncryptedSecret, int64, int) (int64, error)
	insertThresholdShareFn    func(context.Context, int64, int64, []byte) error

	insertKeysFn      func(ctx context.Context, pair models.UserKeyPair) (models.UserKeyPair, error)
	selectKeysFn      func(ctx context.Context, userID int64) (models.UserKeyPair, error)
	hasKeysFn         func(ctx context.Context, userID int64) (bool, error)
	updatePublicKeyFn func(ctx context.Context, userID int64, publicKey []byte) error

	insertPasswordFn func(ctx context.Context, password models.HashedPassword) (models.HashedPassword, error)
	selectPasswordFn func(ctx context.Context, passwordID int64) (models.HashedPassword, error)

	withTxFn func(ctx context.Context, fn func(db.Database) error) error
}

func (m *mockDatabase) UpdatePublicKey(ctx context.Context, userID int64, publicKey []byte) error {
	if m.updatePublicKeyFn != nil {
		return m.updatePublicKeyFn(ctx, userID, publicKey)
	}
	return nil
}

func (m *mockDatabase) SelectPublicKey(ctx context.Context, userID int64) ([]byte, error) {
	if m.selectPublicKeyFn != nil {
		return m.selectPublicKeyFn(ctx, userID)
	}
	return nil, nil
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

func (m *mockDatabase) UpsertSecret(
	ctx context.Context,
	secret models.EncryptedSecret,
	userID int64,
) (models.EncryptedSecret, error) {
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

func (m *mockDatabase) SelectAccessibleSecrets(ctx context.Context, userID int64) ([]models.AccessibleSecret, error) {
	if m.selectAccessibleSecretsFn != nil {
		return m.selectAccessibleSecretsFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockDatabase) SelectSecretForOwner(
	ctx context.Context,
	secretID, userID int64,
) (models.AccessibleSecret, error) {
	if m.selectSecretForOwnerFn != nil {
		return m.selectSecretForOwnerFn(ctx, secretID, userID)
	}
	return models.AccessibleSecret{}, nil
}
func (m *mockDatabase) InsertSecretAccess(ctx context.Context, secretID, userID, grantedBy int64, key []byte) error {
	if m.insertSecretAccessFn != nil {
		return m.insertSecretAccessFn(ctx, secretID, userID, grantedBy, key)
	}
	return nil
}
func (m *mockDatabase) DeleteSecretAccess(context.Context, int64, int64) (int64, error) {
	return 0, nil
}
func (m *mockDatabase) SelectSecretRecipients(context.Context, int64, int64) ([]models.ExistingUser, error) {
	return nil, nil
}
func (m *mockDatabase) UpdateSecretEnvelope(context.Context, int64, int64, []byte) error { return nil }
func (m *mockDatabase) SearchUsers(context.Context, int64, string) ([]models.ExistingUser, error) {
	return nil, nil
}
func (m *mockDatabase) SelectUserBySharingID(ctx context.Context, sharingID string) (models.ExistingUser, error) {
	if m.selectUserBySharingIDFn != nil {
		return m.selectUserBySharingIDFn(ctx, sharingID)
	}
	return models.ExistingUser{}, nil
}

func (m *mockDatabase) InsertThresholdSecret(
	ctx context.Context,
	secret models.EncryptedSecret,
	userID int64,
	threshold int,
) (int64, error) {
	if m.insertThresholdSecretFn != nil {
		return m.insertThresholdSecretFn(ctx, secret, userID, threshold)
	}
	return 1, nil
}
func (m *mockDatabase) InsertThresholdShare(ctx context.Context, secretID, userID int64, share []byte) error {
	if m.insertThresholdShareFn != nil {
		return m.insertThresholdShareFn(ctx, secretID, userID, share)
	}
	return nil
}
func (m *mockDatabase) SelectThresholdSecrets(context.Context, int64) ([]models.ThresholdSecret, error) {
	return nil, nil
}

func (m *mockDatabase) SelectThresholdSecretForParticipant(
	context.Context,
	int64,
	int64,
) (models.ThresholdSecret, error) {
	return models.ThresholdSecret{}, nil
}
func (m *mockDatabase) InsertUnlockRequest(context.Context, int64, int64) (int64, error) {
	return 1, nil
}
func (m *mockDatabase) SelectPendingUnlockRequests(context.Context, int64) ([]models.UnlockRequest, error) {
	return nil, nil
}
func (m *mockDatabase) SelectUnlockRequest(context.Context, int64, int64) (models.ThresholdSecret, error) {
	return models.ThresholdSecret{}, nil
}
func (m *mockDatabase) InsertUnlockContribution(context.Context, int64, int64, []byte) error {
	return nil
}
func (m *mockDatabase) SelectUnlockContributions(context.Context, int64) ([][]byte, error) {
	return nil, nil
}
func (m *mockDatabase) CompleteUnlockRequest(context.Context, int64) error { return nil }

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

func (m *mockDatabase) InsertPassword(
	ctx context.Context,
	password models.HashedPassword,
) (models.HashedPassword, error) {
	if m.insertPasswordFn != nil {
		return m.insertPasswordFn(ctx, password)
	}
	return password, nil
}

func (m *mockDatabase) SelectPassword(
	ctx context.Context,
	passwordID int64,
) (models.HashedPassword, error) {
	if m.selectPasswordFn != nil {
		return m.selectPasswordFn(ctx, passwordID)
	}
	return models.HashedPassword{}, nil
}
