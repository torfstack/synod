package fromdb

import (
	"crypto/x509"
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

func KeyPair(in sqlc.Key) (models.UserKeyPair, error) {
	pub, err := x509.ParsePKCS1PublicKey(in.Public)
	if err != nil {
		return models.UserKeyPair{}, err
	}
	priv, err := x509.ParsePKCS1PrivateKey(in.Private)
	if err != nil {
		return models.UserKeyPair{}, err
	}
	return models.UserKeyPair{
		ID:     &in.ID,
		UserID: in.UserID,
		KeyPair: models.KeyPair{
			Public:  *pub,
			Private: *priv,
		},
	}, nil
}

func tagsSlice(tags string) []string {
	if tags == "" {
		return []string{}
	}
	return strings.Split(tags, ",")
}
