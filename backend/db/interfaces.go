package db

import (
	"context"

	"github.com/torfstack/synod/backend/models"
)

type Database interface {
	WithTx(ctx context.Context) (Database, Transaction)

	DoesUserExist(ctx context.Context, username string) (bool, error)
	InsertUser(ctx context.Context, params models.User) error
	SelectUserByName(ctx context.Context, username string) (models.User, error)
	UpsertSecret(ctx context.Context, secret models.Secret, userID int64) error
	SelectSecrets(ctx context.Context, userID int64) ([]models.Secret, error)
}

type Transaction interface {
	Commit(ctx context.Context)

	// Rollback is a no-op if the transaction has already been committed
	Rollback(ctx context.Context)
}
