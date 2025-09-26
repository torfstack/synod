package db

import (
	"context"
	"crypto/rsa"

	"github.com/torfstack/synod/backend/models"
)

type Database interface {
	WithTx(ctx context.Context, withTx func(Database) error) error

	DoesUserExist(ctx context.Context, username string) (bool, error)
	InsertUser(ctx context.Context, user models.User) (models.ExistingUser, error)
	SelectUserByName(ctx context.Context, username string) (models.ExistingUser, error)
	UpsertSecret(ctx context.Context, secret models.EncryptedSecret, userID int64) (models.EncryptedSecret, error)
	SelectSecrets(ctx context.Context, userID int64) ([]models.EncryptedSecret, error)
	InsertKeys(ctx context.Context, pair models.UserKeyPair) (models.UserKeyPair, error)
	SelectPublicKey(ctx context.Context, userID int64) (rsa.PublicKey, error)
	SelectPrivateKey(ctx context.Context, userID int64) (rsa.PrivateKey, error)
}

type Transaction interface {
	Commit(ctx context.Context)

	// Rollback is a no-op if the transaction has already been committed
	Rollback(ctx context.Context)
}
