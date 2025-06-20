package domain

import (
	"context"

	"github.com/torfstack/kayvault/backend/models"
)

type Service interface {
	UserService
	SecretService
	SessionService
}

type UserService interface {
	DoesUserExist(ctx context.Context, username string) (bool, error)
	InsertUser(ctx context.Context, user models.User) error
	GetUserFromToken(ctx context.Context, token string) (models.ExistingUser, error)
}

type SecretService interface {
	GetSecrets(ctx context.Context, userID int64) ([]models.Secret, error)
	UpsertSecret(ctx context.Context, secret models.Secret, userID int64) error
}

type SessionService interface {
	CreateSession(userID int64) (Session, error)
	GetSession(token string) (*Session, error)
	DeleteSession(token string) error
}
