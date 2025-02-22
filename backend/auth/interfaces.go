package auth

import (
	"context"
	"github.com/torfstack/kayvault/backend/models"
)

type Auth interface {
	GetUser(ctx context.Context, token string) (*models.User, error)
}
