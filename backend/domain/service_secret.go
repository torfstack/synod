package domain

import (
	"context"
	"errors"

	"github.com/torfstack/synod/backend/logging"
	"github.com/torfstack/synod/backend/models"
)

var _ SecretService = &service{}

func (s *service) GetSecrets(ctx context.Context, userId int64) ([]models.Secret, error) {
	return s.database.SelectSecrets(ctx, userId)
}

func (s *service) UpsertSecret(ctx context.Context, secret models.Secret, userID int64) error {
	err := s.database.UpsertSecret(ctx, secret, userID)
	if err != nil {
		logging.Errorf(ctx, "could not upsert secret: %v", err)
		return errors.New("could not upsert secret")
	}
	return nil
}
