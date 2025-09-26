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
}

type UserService interface {
	DoesUserExist(ctx context.Context, username string) (bool, error)
	InsertUser(ctx context.Context, user models.User) (models.ExistingUser, error)
	GetUserFromToken(ctx context.Context, token string) (models.ExistingUser, error)
}

type SecretService interface {
	GetSecrets(ctx context.Context, userID int64) ([]models.Secret, error)
	UpsertSecret(ctx context.Context, secret models.Secret, userID int64) (models.EncryptedSecret, error)
}

type SessionService interface {
	CreateSession(userID int64) (Session, error)
	GetSession(token string) (*Session, error)
	DeleteSession(token string) error
}

type CryptoService interface {
	GenerateKeyPair() (models.KeyPair, error)
	GetPublicKey(ctx context.Context, userID int64) (rsa.PublicKey, error)
	GetPrivateKey(ctx context.Context, userID int64) (rsa.PrivateKey, error)
	EncryptSecret(ctx context.Context, secret models.Secret, key rsa.PublicKey) (models.EncryptedSecret, error)
	DecryptSecret(ctx context.Context, secret models.EncryptedSecret, key rsa.PrivateKey) (models.Secret, error)
}
