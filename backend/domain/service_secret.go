package domain

import (
	"context"

	"github.com/torfstack/kayvault/backend/models"
)

var _ SecretService = &service{}

func (s *service) GetSecrets(ctx context.Context, userId int64) ([]models.Secret, error) {
	return s.database.SelectSecrets(ctx, userId)
}

func (s *service) UpsertSecret(ctx context.Context, secret models.Secret, userID int64) error {
	return s.database.UpsertSecret(ctx, secret, userID)
}
