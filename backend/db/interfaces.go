package db

import (
	"context"

	"github.com/torfstack/synod/backend/models"
)

type Database interface {
	WithTx(ctx context.Context, withTx func(Database) error) error

	DoesUserExist(ctx context.Context, username string) (bool, error)
	InsertUser(ctx context.Context, user models.User) (models.ExistingUser, error)
	SelectUserByName(ctx context.Context, username string) (models.ExistingUser, error)

	UpsertSecret(ctx context.Context, secret models.EncryptedSecret, userID int64) (models.EncryptedSecret, error)
	SelectSecrets(ctx context.Context, userID int64) ([]models.EncryptedSecret, error)
	SelectAccessibleSecrets(ctx context.Context, userID int64) ([]models.AccessibleSecret, error)
	SelectSecretForOwner(ctx context.Context, secretID, userID int64) (models.AccessibleSecret, error)
	InsertSecretAccess(ctx context.Context, secretID, userID, grantedBy int64, encryptedDataKey []byte) error
	DeleteSecretAccess(ctx context.Context, secretID, userID int64) (int64, error)
	SelectSecretRecipients(ctx context.Context, secretID, ownerID int64) ([]models.ExistingUser, error)
	UpdateSecretEnvelope(ctx context.Context, secretID, userID int64, payload []byte) error
	SearchUsers(ctx context.Context, userID int64, search string) ([]models.ExistingUser, error)
	SelectUserBySharingID(ctx context.Context, sharingID string) (models.ExistingUser, error)
	InsertThresholdSecret(
		ctx context.Context,
		secret models.EncryptedSecret,
		userID int64,
		threshold int,
	) (int64, error)
	InsertThresholdShare(ctx context.Context, secretID, userID int64, encryptedShare []byte) error
	SelectThresholdSecrets(ctx context.Context, userID int64) ([]models.ThresholdSecret, error)
	SelectThresholdSecretForParticipant(ctx context.Context, secretID, userID int64) (models.ThresholdSecret, error)
	InsertUnlockRequest(ctx context.Context, secretID, userID int64) (int64, error)
	SelectPendingUnlockRequests(ctx context.Context, userID int64) ([]models.UnlockRequest, error)
	SelectUnlockRequest(ctx context.Context, requestID, userID int64) (models.ThresholdSecret, error)
	InsertUnlockContribution(ctx context.Context, requestID, userID int64, share []byte) error
	SelectUnlockContributions(ctx context.Context, requestID int64) ([][]byte, error)
	CompleteUnlockRequest(ctx context.Context, requestID int64) error

	InsertKeys(ctx context.Context, pair models.UserKeyPair) (models.UserKeyPair, error)
	SelectKeys(ctx context.Context, userID int64) (models.UserKeyPair, error)
	UpdatePublicKey(ctx context.Context, userID int64, publicKey []byte) error
	SelectPublicKey(ctx context.Context, userID int64) ([]byte, error)
	HasKeys(ctx context.Context, userID int64) (bool, error)

	InsertPassword(ctx context.Context, password models.HashedPassword) (models.HashedPassword, error)
	SelectPassword(ctx context.Context, passwordID int64) (models.HashedPassword, error)
}

type Transaction interface {
	Commit(ctx context.Context)

	// Rollback is a no-op if the transaction has already been committed
	Rollback(ctx context.Context)
}
