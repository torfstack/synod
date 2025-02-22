package todb

import (
	"github.com/torfstack/kayvault/backend/models"
	sqlc "github.com/torfstack/kayvault/sql/gen"
)

func Secret(in models.Secret) sqlc.Secret {
	return sqlc.Secret{
		ID:    in.ID,
		Value: []byte(in.Value),
		Key:   in.Key,
		Url:   in.Url,
		Tags:  tagsString(in.Tags),
	}
}

func InsertSecretParams(in models.Secret, userID int32) sqlc.InsertSecretParams {
	return sqlc.InsertSecretParams{
		Value:  []byte(in.Value),
		Key:    in.Key,
		Url:    in.Url,
		Tags:   tagsString(in.Tags),
		UserID: userID,
	}
}

func UpdateSecretParams(in models.Secret, userID int32) sqlc.UpdateSecretParams {
	return sqlc.UpdateSecretParams{
		ID:     in.ID,
		Value:  []byte(in.Value),
		Key:    in.Key,
		Url:    in.Url,
		Tags:   tagsString(in.Tags),
		UserID: userID,
	}
}

func tagsString(tags []string) string {
	t := ""
	for _, tag := range tags {
		t += tag + ","
	}
	return t[:len(t)-1]
}
