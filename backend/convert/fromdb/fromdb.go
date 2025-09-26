package fromdb

import (
	"strings"

	"github.com/torfstack/synod/backend/models"
	sqlc "github.com/torfstack/synod/sql/gen"
)

func Secret(in sqlc.Secret) models.Secret {
	return models.Secret{
		ID:    &in.ID,
		Value: string(in.Value),
		Key:   in.Key,
		Url:   in.Url,
		Tags:  tagsSlice(in.Tags),
	}
}

func Secrets(in []sqlc.Secret) models.Secrets {
	out := make([]models.Secret, len(in))
	for i, s := range in {
		out[i] = Secret(s)
	}
	return out
}

func User(in sqlc.User) models.ExistingUser {
	return models.ExistingUser{
		ID: in.ID,
		User: models.User{
			Subject:  in.Subject,
			Email:    in.Email,
			FullName: in.FullName,
		},
	}
}

func KeyPair(in sqlc.Key) models.UserKeyPair {
	return models.UserKeyPair{
		ID:     &in.ID,
		UserID: in.UserID,
		KeyPair: models.KeyPair{
			Public:  in.Public,
			Private: in.Private,
		},
	}
}

func tagsSlice(tags string) []string {
	if tags == "" {
		return []string{}
	}
	return strings.Split(tags, ",")
}
