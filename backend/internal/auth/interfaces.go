package auth

import (
	"context"
	"main/internal/models"
)

type Auth interface {
	GetUser(ctx context.Context, token string) (*models.User, error)
}
