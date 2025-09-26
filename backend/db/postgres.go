package db

import (
	"context"
	"crypto/rsa"
	"crypto/x509"

	"github.com/jackc/pgx/v5"
	"github.com/torfstack/synod/backend/convert/fromdb"
	"github.com/torfstack/synod/backend/convert/todb"
	"github.com/torfstack/synod/backend/models"
	sqlc "github.com/torfstack/synod/sql/gen"
)

type database struct {
	connStr string
	conn    *pgx.Conn
	tx      *transaction
}

var _ Database = (*database)(nil)

func NewDatabase(connStr string) Database {
	return &database{connStr: connStr}
}

func (d *database) WithTx(ctx context.Context, withTx func(Database) error) error {
	if d.tx != nil {
		return withTx(d)
	}

	var conn *pgx.Conn
	if d.conn != nil {
		conn = d.conn
	} else {
		var err error
		conn, err = pgx.Connect(ctx, d.connStr)
		if err != nil {
			return err
		}
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	trans := &transaction{conn: conn, tx: tx}
	defer tx.Rollback(ctx)
	err = withTx(&database{connStr: d.connStr, conn: conn, tx: trans})
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (d *database) CommitTransaction(ctx context.Context) error {
	if d.tx == nil {
		return nil
	}
	defer func(context.Context) {
		d.tx.Rollback(ctx)
		_ = (*d.conn).Close(ctx)
	}(ctx)
	d.tx.Commit(ctx)
	return nil
}

func (d *database) DoesUserExist(ctx context.Context, username string) (bool, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return false, err
	}
	return q.DoesUserExist(ctx, username)
}

func (d *database) InsertUser(ctx context.Context, user models.User) (models.ExistingUser, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return models.ExistingUser{}, err
	}
	params := todb.InsertUserParams(user)
	dbUser, err := q.InsertUser(ctx, params)
	return fromdb.User(dbUser), err
}

func (d *database) SelectUserByName(ctx context.Context, username string) (models.ExistingUser, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
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
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
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
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return []models.EncryptedSecret{}, err
	}
	dbSecrets, err := q.SelectSecrets(ctx, userID)
	return fromdb.Secrets(dbSecrets), err
}

func (d *database) InsertKeys(ctx context.Context, pair models.UserKeyPair) (models.UserKeyPair, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return models.UserKeyPair{}, err
	}
	params := todb.InsertKeysParams(pair)
	dbKeys, err := q.InsertKeys(ctx, params)
	if err != nil {
		return models.UserKeyPair{}, err
	}
	return fromdb.KeyPair(dbKeys)
}

func (d *database) SelectPublicKey(ctx context.Context, userID int64) (rsa.PublicKey, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return rsa.PublicKey{}, err
	}
	dbKey, err := q.SelectPublicKeyForUser(ctx, userID)
	if err != nil {
		return rsa.PublicKey{}, err
	}
	parsedKey, err := x509.ParsePKCS1PublicKey(dbKey)
	if err != nil {
		return rsa.PublicKey{}, err
	}
	return *parsedKey, nil
}

func (d *database) SelectPrivateKey(ctx context.Context, userID int64) (rsa.PrivateKey, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return rsa.PrivateKey{}, err
	}
	dbKey, err := q.SelectPrivateKeyForUser(ctx, userID)
	if err != nil {
		return rsa.PrivateKey{}, err
	}
	parsedKey, err := x509.ParsePKCS1PrivateKey(dbKey)
	if err != nil {
		return rsa.PrivateKey{}, err
	}
	return *parsedKey, nil
}

func startQuery(ctx context.Context, d *database) (*sqlc.Queries, error) {
	if d.tx != nil {
		return sqlc.New(d.tx.SqlTx()), nil
	}
	if d.conn == nil {
		conn, err := pgx.Connect(ctx, d.connStr)
		if err != nil {
			return nil, err
		}
		d.conn = conn
	}
	return sqlc.New(d.conn), nil
}

func endQuery(ctx context.Context, d *database) {
	if d.conn != nil && d.tx == nil {
		_ = (*d.conn).Close(ctx)
		d.conn = nil
	}
}
