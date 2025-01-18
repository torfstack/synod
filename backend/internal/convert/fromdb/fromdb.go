package fromdb

import (
	"github.com/torfstack/kayvault/internal/models"
	sqlc "github.com/torfstack/kayvault/sql/gen"
)

func Secret(in sqlc.Secret) models.Secret {
	return models.Secret{
		ID:    in.ID,
		Value: string(in.Value),
		Key:   in.Key,
		Url:   in.Url,
	}
}

func Secrets(in []sqlc.Secret) []models.Secret {
	out := make([]models.Secret, len(in))
	for i, s := range in {
		out[i] = Secret(s)
	}
	return out
}

func User(in sqlc.User) models.User {
	return models.User{
		ID:       in.ID,
		Username: in.Username,
	}
}
