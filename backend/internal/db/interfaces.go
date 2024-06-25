package db

import (
	"context"
	sqlc "main/sql/gen"
)

type Database interface {
	Connect(ctx context.Context) (Connection, error)
}

type Connection interface {
	Close(ctx context.Context) error
	Queries() Queries
}

type Queries interface {
	SelectSecrets(ctx context.Context, userId int32) ([]sqlc.Secret, error)
	InsertSecret(ctx context.Context, arg sqlc.InsertSecretParams) error

	SelectUserByName(ctx context.Context, username string) (sqlc.User, error)
	InsertUser(ctx context.Context, username string) error
	DoesUserExist(ctx context.Context, username string) (bool, error)
}
