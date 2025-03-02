package domain

import (
	"context"

	"github.com/torfstack/kayvault/backend/db"
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
	GetUserFromToken(ctx context.Context, token string) (*models.User, error)
}

type SecretService interface {
	GetSecrets(ctx context.Context, userId int32) ([]models.Secret, error)
	UpsertSecret(ctx context.Context, secret models.Secret, userId int32) error
}

type SessionService interface {
	CreateSession(user int32) (*Session, error)
	GetSession(token string) (*Session, error)
	DeleteSession(token string) error
}

type service struct {
	database db.Database
	sessions sessionStore
}

func NewDomainService(db db.Database) Service {
	return &service{
		database: db,
		sessions: make(sessionStore),
	}
}
