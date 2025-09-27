package domain

import (
	"context"
	"crypto/rsa"

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
	GetUserFromToken(ctx context.Context, token string) (models.ExistingUser, error)
}

type SecretService interface {
	GetSecrets(ctx context.Context, userID int64, key *rsa.PrivateKey) ([]models.Secret, error)
	UpsertSecret(ctx context.Context, secret models.Secret, userID int64, key *rsa.PrivateKey) (models.EncryptedSecret, error)
}

type SessionService interface {
	CreateSession(ctx context.Context, userID int64) (Session, error)
	GetSession(token string) (*Session, error)
	DeleteSession(token string) error
}

type CryptoService interface {
	GenerateKeyPair() (models.KeyPair, error)
	EncryptSecret(ctx context.Context, secret models.Secret, key *rsa.PublicKey) (models.EncryptedSecret, error)
	DecryptSecret(ctx context.Context, secret models.EncryptedSecret, key *rsa.PrivateKey) (models.Secret, error)
}

type SetupService interface {
	IsUserSetup(ctx context.Context, session Session) (bool, error)
	SetupUserPlain(ctx context.Context, session Session) error
	SetupUserWithPassword(ctx context.Context, session Session, password string) error
	UnsealWithPassword(ctx context.Context, session *Session, password string) error
}
