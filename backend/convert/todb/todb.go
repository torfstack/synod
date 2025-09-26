package todb

import (
	"crypto/x509"

	"github.com/torfstack/synod/backend/models"
	sqlc "github.com/torfstack/synod/sql/gen"
)

func Secret(in models.Secret) sqlc.Secret {
	return sqlc.Secret{
		ID:    *in.ID,
		Value: []byte(in.Value),
		Key:   in.Key,
		Url:   in.Url,
		Tags:  tagsString(in.Tags),
	}
}

func InsertSecretParams(in models.EncryptedSecret, userID int64) sqlc.InsertSecretParams {
	return sqlc.InsertSecretParams{
		Value:  []byte(in.Value),
		Key:    in.Key,
		Url:    in.Url,
		Tags:   tagsString(in.Tags),
		UserID: userID,
	}
}

func UpdateSecretParams(in models.EncryptedSecret, userID int64) sqlc.UpdateSecretParams {
	return sqlc.UpdateSecretParams{
		ID:     *in.ID,
		Value:  []byte(in.Value),
		Key:    in.Key,
		Url:    in.Url,
		Tags:   tagsString(in.Tags),
		UserID: userID,
	}
}

func InsertUserParams(in models.User) sqlc.InsertUserParams {
	return sqlc.InsertUserParams{
		Subject:  in.Subject,
		Email:    in.Email,
		FullName: in.FullName,
	}
}

func InsertKeysParams(in models.UserKeyPair) sqlc.InsertKeysParams {
	pub := x509.MarshalPKCS1PublicKey(&in.Public)
	priv := x509.MarshalPKCS1PrivateKey(&in.Private)
	return sqlc.InsertKeysParams{
		UserID:  in.UserID,
		Public:  pub,
		Private: priv,
	}
}

func tagsString(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	t := ""
	for _, tag := range tags {
		t += tag + ","
	}
	return t[:len(t)-1]
}
