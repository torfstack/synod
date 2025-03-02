package domain

import (
	"context"

	"github.com/torfstack/kayvault/backend/convert/fromdb"
	"github.com/torfstack/kayvault/backend/convert/todb"
	"github.com/torfstack/kayvault/backend/models"
)

var _ SecretService = &service{}

func (s *service) GetSecrets(ctx context.Context, userId int32) ([]models.Secret, error) {
	dbSecrets, err := s.database.SelectSecrets(ctx, userId)
	if err != nil {
		return nil, err
	}
	return fromdb.Secrets(dbSecrets), nil
}

func (s *service) UpsertSecret(ctx context.Context, secret models.Secret, userId int32) error {
	if secret.ID != 0 {
		return s.database.UpdateSecret(ctx, todb.UpdateSecretParams(secret, userId))
	}
	return s.database.InsertSecret(ctx, todb.InsertSecretParams(secret, userId))
}
