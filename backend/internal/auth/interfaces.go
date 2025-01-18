package auth

import (
	"context"
	"github.com/torfstack/kayvault/internal/models"
)

type Auth interface {
	GetUser(ctx context.Context, token string) (*models.User, error)
}
