package domain

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/models"
)

type Service interface {
	UserService
	SecretService
	SessionService
	SetupService
}

type UserService interface {
	DoesUserExist(ctx context.Context, username string) (bool, error)
	InsertUser(ctx context.Context, user models.User) (models.ExistingUser, error)
	GetUserFromToken(ctx context.Context, token *oidc.IDToken) (models.ExistingUser, error)
}

type SecretService interface {
	GetSecrets(ctx context.Context, userID int64, cipher *crypto.AsymmetricCipher) ([]models.Secret, error)
	UpsertSecret(
		ctx context.Context,
		secret models.Secret,
		userID int64,
		cipher *crypto.AsymmetricCipher,
	) (models.EncryptedSecret, error)
}

type SessionService interface {
	CreateSession(ctx context.Context, userID int64) (Session, error)
	GetSession(token string) (*Session, error)
	DeleteSession(token string) error
}

type CryptoService interface {
	EncryptSecret(
		ctx context.Context,
		secret models.Secret,
		cipher *crypto.AsymmetricCipher,
	) (models.EncryptedSecret, error)
	DecryptSecret(
		ctx context.Context,
		secret models.EncryptedSecret,
		cipher *crypto.AsymmetricCipher,
	) (models.Secret, error)
}

type SetupService interface {
	IsUserSetup(ctx context.Context, session Session) (bool, error)
	SetupUserPlain(ctx context.Context, session Session) error
	SetupUserWithPassword(ctx context.Context, session Session, password string) error
	UnsealWithPassword(ctx context.Context, session *Session, password string) error
}
