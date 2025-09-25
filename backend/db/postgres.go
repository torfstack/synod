package db

import (
	"context"

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

func (d *database) WithTx(ctx context.Context) (Database, Transaction) {
	if d.tx != nil {
		return d, d.tx
	}

	var conn *pgx.Conn
	if d.conn != nil {
		conn = d.conn
	} else {
		var err error
		conn, err = pgx.Connect(ctx, d.connStr)
		if err != nil {
			return nil, nil
		}
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, nil
	}
	trans := &transaction{conn: conn, tx: tx}
	return &database{connStr: d.connStr, conn: conn, tx: trans}, trans
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

func (d *database) InsertUser(ctx context.Context, user models.User) error {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return err
	}
	params := todb.InsertUserParams(user)
	return q.InsertUser(ctx, params)
}

func (d *database) SelectUserByName(ctx context.Context, username string) (models.User, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return models.User{}, err
	}
	dbUser, err := q.SelectUserByName(ctx, username)
	return fromdb.User(dbUser), err
}

func (d *database) UpsertSecret(ctx context.Context, secret models.Secret, userID int64) error {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return err
	}
	if secret.ID == nil || *secret.ID == 0 {
		params := todb.InsertSecretParams(secret, userID)
		return q.InsertSecret(ctx, params)
	} else {
		params := todb.UpdateSecretParams(secret, userID)
		return q.UpdateSecret(ctx, params)
	}
}

func (d *database) SelectSecrets(ctx context.Context, userID int64) ([]models.Secret, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return []models.Secret{}, err
	}
	dbSecrets, err := q.SelectSecrets(ctx, userID)
	return fromdb.Secrets(dbSecrets), err
}

func (d *database) InsertKeys(ctx context.Context, pair models.UserKeyPair) error {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return err
	}
	params := todb.InsertKeysParams(pair)
	return q.InsertKeys(ctx, params)
}

func (d *database) SelectPublicKey(ctx context.Context, userID int64) ([]byte, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return []byte{}, err
	}
	return q.SelectPublicKeyForUser(ctx, userID)
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
