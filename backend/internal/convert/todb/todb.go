package todb

import (
	"main/internal/models"
	sqlc "main/sql/gen"
)

func Secret(in models.Secret) sqlc.Secret {
	return sqlc.Secret{
		ID:    in.ID,
		Value: []byte(in.Value),
		Key:   in.Key,
		Url:   in.Url,
	}
}

func InsertSecretParams(in models.Secret, userId int32) sqlc.InsertSecretParams {
	return sqlc.InsertSecretParams{
		Value:  []byte(in.Value),
		Key:    in.Key,
		Url:    in.Url,
		UserID: userId,
	}
}
