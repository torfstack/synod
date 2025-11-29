package fromdb

import (
	"strings"

	"github.com/torfstack/synod/backend/models"
	sqlc "github.com/torfstack/synod/sql/gen"
)

func Secret(in sqlc.Secret) models.EncryptedSecret {
	return models.EncryptedSecret{
		ID:    &in.ID,
		Value: string(in.Value),
		Key:   in.Key,
		Url:   in.Url,
		Tags:  tagsSlice(in.Tags),
	}
}

func Secrets(in []sqlc.Secret) []models.EncryptedSecret {
	out := make([]models.EncryptedSecret, len(in))
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
	userKeyPair := models.UserKeyPair{
		ID:          &in.ID,
		Type:        models.KeyType(in.Type),
		UserID:      in.UserID,
		KeyMaterial: in.KeyMaterial,
	}
	if in.PasswordID.Valid {
		userKeyPair.PasswordID = &in.PasswordID.Int64
	}
	return userKeyPair
}

func HashedPassword(in sqlc.Password) models.HashedPassword {
	return models.HashedPassword{
		ID:         &in.ID,
		Hash:       in.Hash,
		Salt:       in.Salt,
		Iterations: in.Iterations,
	}
}

func tagsSlice(tags string) []string {
	if tags == "" {
		return []string{}
	}
	return strings.Split(tags, ",")
}
