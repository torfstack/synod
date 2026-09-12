package db

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/torfstack/synod/backend/convert/fromdb"
	"github.com/torfstack/synod/backend/convert/todb"
	"github.com/torfstack/synod/backend/models"
	sqlc "github.com/torfstack/synod/sql/gen"
)

type database struct {
	connStr string
	pool    *pgxpool.Pool
	tx      *transaction
}

var _ Database = (*database)(nil)

func NewDatabase(ctx context.Context, connStr string) (Database, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &database{connStr: connStr, pool: pool}, nil
}

func (d *database) WithTx(ctx context.Context, withTx func(Database) error) error {
	if d.tx != nil {
		return withTx(d)
	}

	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}
	trans := &transaction{tx: tx}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)
	if err = withTx(&database{connStr: d.connStr, pool: d.pool, tx: trans}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (d *database) DoesUserExist(ctx context.Context, username string) (bool, error) {
	q, err := startQuery(d)
	if err != nil {
		return false, err
	}
	return q.DoesUserExist(ctx, username)
}

func (d *database) InsertUser(ctx context.Context, user models.User) (models.ExistingUser, error) {
	q, err := startQuery(d)
	if err != nil {
		return models.ExistingUser{}, err
	}
	params := todb.InsertUserParams(user)
	dbUser, err := q.InsertUser(ctx, params)
	return fromdb.User(dbUser), err
}

func (d *database) SelectUserByName(ctx context.Context, username string) (models.ExistingUser, error) {
	q, err := startQuery(d)
	if err != nil {
		return models.ExistingUser{}, err
	}
	dbUser, err := q.SelectUserByName(ctx, username)
	return fromdb.User(dbUser), err
}

func (d *database) UpsertSecret(
	ctx context.Context,
	secret models.EncryptedSecret,
	userID int64,
) (models.EncryptedSecret, error) {
	q, err := startQuery(d)
	if err != nil {
		return models.EncryptedSecret{}, err
	}
	var dbSecret sqlc.Secret
	if secret.ID == nil || *secret.ID == 0 {
		params := todb.InsertSecretParams(secret, userID)
		dbSecret, err = q.InsertSecret(ctx, params)
	} else {
		params := todb.UpdateSecretParams(secret, userID)
		dbSecret, err = q.UpdateSecret(ctx, params)
	}
	return fromdb.Secret(dbSecret), err
}

func (d *database) SelectSecrets(ctx context.Context, userID int64) ([]models.EncryptedSecret, error) {
	q, err := startQuery(d)
	if err != nil {
		return []models.EncryptedSecret{}, err
	}
	dbSecrets, err := q.SelectSecrets(ctx, userID)
	return fromdb.Secrets(dbSecrets), err
}

func (d *database) SelectAccessibleSecrets(ctx context.Context, userID int64) ([]models.AccessibleSecret, error) {
	q, err := startQuery(d)
	if err != nil {
		return nil, err
	}
	rows, err := q.SelectAccessibleSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]models.AccessibleSecret, len(rows))
	for i, row := range rows {
		out[i] = models.AccessibleSecret{
			ID:               row.ID,
			OwnerID:          row.UserID,
			EncryptedPayload: row.Value,
			EncryptedDataKey: row.EncryptedDataKey,
			Envelope:         row.Envelope,
			Owned:            row.Owned,
		}
	}
	return out, nil
}

func (d *database) SelectSecretForOwner(ctx context.Context, secretID, userID int64) (models.AccessibleSecret, error) {
	q, err := startQuery(d)
	if err != nil {
		return models.AccessibleSecret{}, err
	}
	row, err := q.SelectSecretForOwner(ctx, sqlc.SelectSecretForOwnerParams{ID: secretID, UserID: userID})
	if err != nil {
		return models.AccessibleSecret{}, err
	}
	return models.AccessibleSecret{
		ID:               row.ID,
		OwnerID:          row.UserID,
		EncryptedPayload: row.Value,
		EncryptedDataKey: row.EncryptedDataKey,
		Envelope:         row.Envelope,
		Legacy: fromdb.Secret(
			sqlc.Secret{
				ID:            row.ID,
				Value:         row.Value,
				Key:           row.Key,
				Url:           row.Url,
				Tags:          row.Tags,
				UserID:        row.UserID,
				SecretSharing: row.SecretSharing,
				Envelope:      row.Envelope,
				CreatedAt:     row.CreatedAt,
				UpdatedAt:     row.UpdatedAt,
			},
		),
	}, nil
}

func (d *database) InsertSecretAccess(
	ctx context.Context,
	secretID, userID, grantedBy int64,
	encryptedDataKey []byte,
) error {
	q, err := startQuery(d)
	if err != nil {
		return err
	}
	return q.InsertSecretAccess(
		ctx,
		sqlc.InsertSecretAccessParams{
			SecretID:         secretID,
			UserID:           userID,
			EncryptedDataKey: encryptedDataKey,
			GrantedBy:        grantedBy,
		},
	)
}

func (d *database) DeleteSecretAccess(ctx context.Context, secretID, userID int64) (int64, error) {
	q, err := startQuery(d)
	if err != nil {
		return 0, err
	}
	return q.DeleteSecretAccess(ctx, sqlc.DeleteSecretAccessParams{SecretID: secretID, UserID: userID})
}

func (d *database) SelectSecretRecipients(ctx context.Context, secretID, ownerID int64) ([]models.ExistingUser, error) {
	q, err := startQuery(d)
	if err != nil {
		return nil, err
	}
	rows, err := q.SelectSecretRecipients(ctx, sqlc.SelectSecretRecipientsParams{ID: secretID, UserID: ownerID})
	if err != nil {
		return nil, err
	}
	out := make([]models.ExistingUser, len(rows))
	for i, row := range rows {
		out[i] = fromdb.User(row)
	}
	return out, nil
}

func (d *database) UpdateSecretEnvelope(ctx context.Context, secretID, userID int64, payload []byte) error {
	q, err := startQuery(d)
	if err != nil {
		return err
	}
	return q.UpdateSecretEnvelope(ctx, sqlc.UpdateSecretEnvelopeParams{Value: payload, ID: secretID, UserID: userID})
}

func (d *database) SearchUsers(ctx context.Context, userID int64, search string) ([]models.ExistingUser, error) {
	q, err := startQuery(d)
	if err != nil {
		return nil, err
	}
	rows, err := q.SearchUsers(ctx, sqlc.SearchUsersParams{ID: userID, Lower: search})
	if err != nil {
		return nil, err
	}
	out := make([]models.ExistingUser, len(rows))
	for i, row := range rows {
		out[i] = fromdb.User(row)
	}
	return out, nil
}

func (d *database) SelectUserBySharingID(ctx context.Context, sharingID string) (models.ExistingUser, error) {
	id, err := uuid.Parse(sharingID)
	if err != nil {
		return models.ExistingUser{}, err
	}
	q, err := startQuery(d)
	if err != nil {
		return models.ExistingUser{}, err
	}
	user, err := q.SelectUserBySharingID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	return fromdb.User(user), err
}

func (d *database) InsertThresholdSecret(
	ctx context.Context,
	secret models.EncryptedSecret,
	userID int64,
	threshold int,
) (int64, error) {
	q, err := startQuery(d)
	if err != nil {
		return 0, err
	}
	stored, err := q.InsertThresholdSecret(ctx, sqlc.InsertThresholdSecretParams{
		Value: []byte(
			secret.Value,
		),
		Key:           secret.Key,
		Url:           secret.Url,
		Tags:          strings.Join(secret.Tags, ","),
		UserID:        userID,
		SecretSharing: pgtype.Int4{Int32: int32(threshold), Valid: true},
	})
	return stored.ID, err
}

func (d *database) InsertThresholdShare(ctx context.Context, secretID, userID int64, encryptedShare []byte) error {
	q, err := startQuery(d)
	if err != nil {
		return err
	}
	return q.InsertThresholdShare(
		ctx,
		sqlc.InsertThresholdShareParams{SecretID: secretID, UserID: userID, EncryptedShare: encryptedShare},
	)
}

func (d *database) SelectThresholdSecrets(ctx context.Context, userID int64) ([]models.ThresholdSecret, error) {
	q, err := startQuery(d)
	if err != nil {
		return nil, err
	}
	rows, err := q.SelectThresholdSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]models.ThresholdSecret, len(rows))
	for i, row := range rows {
		out[i] = models.ThresholdSecret{
			ID:        row.ID,
			OwnerID:   row.UserID,
			Threshold: int(row.SecretSharing.Int32),
			Key:       row.Key,
			Url:       row.Url,
			Tags:      splitTags(row.Tags),
		}
	}
	return out, nil
}

func (d *database) SelectThresholdSecretForParticipant(
	ctx context.Context,
	secretID, userID int64,
) (models.ThresholdSecret, error) {
	q, err := startQuery(d)
	if err != nil {
		return models.ThresholdSecret{}, err
	}
	row, err := q.SelectThresholdSecretForParticipant(
		ctx,
		sqlc.SelectThresholdSecretForParticipantParams{ID: secretID, UserID: userID},
	)
	if err != nil {
		return models.ThresholdSecret{}, err
	}
	return models.ThresholdSecret{
		ID:               row.ID,
		OwnerID:          row.UserID,
		EncryptedPayload: row.Value,
		EncryptedShare:   row.EncryptedShare,
		Threshold:        int(row.SecretSharing.Int32),
		Key:              row.Key,
		Url:              row.Url,
		Tags:             splitTags(row.Tags),
	}, nil
}

func (d *database) InsertUnlockRequest(ctx context.Context, secretID, userID int64) (int64, error) {
	q, err := startQuery(d)
	if err != nil {
		return 0, err
	}
	row, err := q.InsertUnlockRequest(ctx, sqlc.InsertUnlockRequestParams{ID: secretID, RequestedBy: userID})
	return row.ID, err
}

func (d *database) SelectPendingUnlockRequests(ctx context.Context, userID int64) ([]models.UnlockRequest, error) {
	q, err := startQuery(d)
	if err != nil {
		return nil, err
	}
	rows, err := q.SelectPendingUnlockRequests(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]models.UnlockRequest, len(rows))
	for i, row := range rows {
		out[i] = models.UnlockRequest{
			ID:            row.ID,
			SecretID:      row.SecretID,
			RequesterID:   row.RequestedBy,
			RequesterName: row.RequesterName,
			SecretName:    row.Key,
			Threshold:     int(row.SecretSharing.Int32),
			Contributions: int(row.Contributions),
			Contributed:   row.Contributed,
			ExpiresAt:     row.ExpiresAt.Time,
		}
	}
	return out, nil
}

func (d *database) SelectUnlockRequest(ctx context.Context, requestID, userID int64) (models.ThresholdSecret, error) {
	q, err := startQuery(d)
	if err != nil {
		return models.ThresholdSecret{}, err
	}
	row, err := q.SelectUnlockRequest(ctx, sqlc.SelectUnlockRequestParams{ID: requestID, RequestedBy: userID})
	if err != nil {
		return models.ThresholdSecret{}, err
	}
	return models.ThresholdSecret{
		ID:               row.SecretID,
		EncryptedPayload: row.EncryptedPayload,
		Threshold:        int(row.SecretSharing.Int32),
	}, nil
}

func (d *database) InsertUnlockContribution(ctx context.Context, requestID, userID int64, share []byte) error {
	q, err := startQuery(d)
	if err != nil {
		return err
	}
	return q.InsertUnlockContribution(
		ctx,
		sqlc.InsertUnlockContributionParams{ID: requestID, UserID: userID, Share: share},
	)
}

func (d *database) SelectUnlockContributions(ctx context.Context, requestID int64) ([][]byte, error) {
	q, err := startQuery(d)
	if err != nil {
		return nil, err
	}
	return q.SelectUnlockContributions(ctx, requestID)
}

func (d *database) CompleteUnlockRequest(ctx context.Context, requestID int64) error {
	q, err := startQuery(d)
	if err != nil {
		return err
	}
	return q.CompleteUnlockRequest(ctx, requestID)
}

func splitTags(value string) []string {
	if value == "" {
		return []string{}
	}
	return strings.Split(value, ",")
}

func (d *database) InsertKeys(ctx context.Context, pair models.UserKeyPair) (models.UserKeyPair, error) {
	q, err := startQuery(d)
	if err != nil {
		return models.UserKeyPair{}, err
	}
	params := todb.InsertKeysParams(pair)
	dbKeys, err := q.InsertKeys(ctx, params)
	if err != nil {
		return models.UserKeyPair{}, err
	}
	return fromdb.KeyPair(dbKeys), nil
}

func (d *database) SelectKeys(ctx context.Context, userID int64) (models.UserKeyPair, error) {
	q, err := startQuery(d)
	if err != nil {
		return models.UserKeyPair{}, err
	}
	dbKeys, err := q.SelectKeys(ctx, userID)
	if err != nil {
		return models.UserKeyPair{}, err
	}
	return fromdb.KeyPair(dbKeys), nil
}

func (d *database) HasKeys(ctx context.Context, userID int64) (bool, error) {
	q, err := startQuery(d)
	if err != nil {
		return false, err
	}
	return q.HasKeys(ctx, userID)
}

func (d *database) UpdatePublicKey(ctx context.Context, userID int64, publicKey []byte) error {
	q, err := startQuery(d)
	if err != nil {
		return err
	}
	return q.UpdatePublicKey(ctx, sqlc.UpdatePublicKeyParams{UserID: userID, PublicKey: publicKey})
}

func (d *database) SelectPublicKey(ctx context.Context, userID int64) ([]byte, error) {
	q, err := startQuery(d)
	if err != nil {
		return nil, err
	}
	key, err := q.SelectPublicKey(ctx, userID)
	return key.PublicKey, err
}

func (d *database) InsertPassword(ctx context.Context, password models.HashedPassword) (models.HashedPassword, error) {
	q, err := startQuery(d)
	if err != nil {
		return models.HashedPassword{}, err
	}
	params := todb.InsertPasswordParams(password)
	dbPassword, err := q.InsertPassword(ctx, params)
	if err != nil {
		return models.HashedPassword{}, err
	}
	return fromdb.HashedPassword(dbPassword), nil
}

func (d *database) SelectPassword(ctx context.Context, passwordID int64) (models.HashedPassword, error) {
	q, err := startQuery(d)
	if err != nil {
		return models.HashedPassword{}, err
	}
	dbPassword, err := q.SelectPassword(ctx, passwordID)
	if err != nil {
		return models.HashedPassword{}, err
	}
	return fromdb.HashedPassword(dbPassword), nil
}

func startQuery(d *database) (*sqlc.Queries, error) {
	if d.tx != nil {
		return sqlc.New(d.tx.SqlTx()), nil
	}
	return sqlc.New(d.pool), nil
}
