package domain

import (
	"context"
	"crypto/rsa"
	"errors"

	"github.com/torfstack/synod/backend/logging"
	"github.com/torfstack/synod/backend/models"
)

var _ SecretService = &service{}

func (s *service) GetSecrets(ctx context.Context, userID int64, key *rsa.PrivateKey) ([]models.Secret, error) {
	encryptedSecrets, err := s.database.SelectSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return []models.Secret{}, errors.New("need private key to decrypt secrets")
	}

	secrets := make([]models.Secret, len(encryptedSecrets))
	for i, encryptedSecret := range encryptedSecrets {
		decrypted, err := s.DecryptSecret(ctx, encryptedSecret, key)
		if err != nil {
			return nil, err
		}
		secrets[i] = decrypted
	}
	return secrets, nil
}

func (s *service) UpsertSecret(ctx context.Context, secret models.Secret, userID int64, key *rsa.PrivateKey) (models.EncryptedSecret, error) {
	if key == nil {
		return models.EncryptedSecret{}, errors.New("need public key to encrypt secrets")
	}
	pub := key.PublicKey
	encryptedSecret, err := s.EncryptSecret(ctx, secret, &pub)
	if err != nil {
		return models.EncryptedSecret{}, err
	}
	upsertSecret, err := s.database.UpsertSecret(ctx, encryptedSecret, userID)
	if err != nil {
		logging.Errorf(ctx, "could not upsert secret: %v", err)
		return models.EncryptedSecret{}, errors.New("could not upsert secret")
	}
	return upsertSecret, nil
}
