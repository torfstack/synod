package todb

import (
	"github.com/jackc/pgx/v5/pgtype"
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
	params := sqlc.InsertKeysParams{
		UserID:  in.UserID,
		Type:    int32(in.Type),
		Public:  in.Public,
		Private: in.Private,
	}
	if in.PasswordID != nil {
		params.PasswordID = pgtype.Int8{
			Int64: *in.PasswordID,
			Valid: true,
		}
	}
	return params
}

func InsertPasswordParams(in models.HashedPassword) sqlc.InsertPasswordParams {
	return sqlc.InsertPasswordParams{
		Hash:       in.Hash,
		Salt:       in.Salt,
		Iterations: in.Iterations,
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
