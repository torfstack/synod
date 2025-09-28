package db

import (
	"context"

	"github.com/jackc/pgx/v5"
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
	err = withTx(&database{connStr: d.connStr, pool: d.pool, tx: trans})
	if err != nil {
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
