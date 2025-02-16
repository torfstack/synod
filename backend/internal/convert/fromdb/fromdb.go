package fromdb

import (
	"github.com/torfstack/kayvault/internal/models"
	sqlc "github.com/torfstack/kayvault/sql/gen"
	"strings"
)

func Secret(in sqlc.Secret) models.Secret {
	return models.Secret{
		ID:    in.ID,
		Value: string(in.Value),
		Key:   in.Key,
		Url:   in.Url,
		Tags:  tagsSlice(in.Tags),
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
		Subject:  in.Subject,
		Email:    in.Email,
		FullName: in.FullName,
	}
}

func tagsSlice(tags string) []string {
	if tags == "" {
		return []string{}
	}
	return strings.Split(tags, ",")
}
