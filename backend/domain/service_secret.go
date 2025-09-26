package domain

import (
	"context"
	"errors"

	"github.com/torfstack/synod/backend/logging"
	"github.com/torfstack/synod/backend/models"
)

var _ SecretService = &service{}

func (s *service) GetSecrets(ctx context.Context, userID int64) ([]models.Secret, error) {
	encryptedSecrets, err := s.database.SelectSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}

	key, err := s.database.SelectPrivateKey(ctx, userID)
	if err != nil {
		return nil, err
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

func (s *service) UpsertSecret(ctx context.Context, secret models.Secret, userID int64) (models.EncryptedSecret, error) {
	key, err := s.database.SelectPublicKey(ctx, userID)
	if err != nil {
		return models.EncryptedSecret{}, err
	}
	encryptedSecret, err := s.EncryptSecret(ctx, secret, key)
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
