package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/torfstack/kayvault/backend/convert/fromdb"
	"github.com/torfstack/kayvault/backend/convert/todb"
	"github.com/torfstack/kayvault/backend/models"
	sqlc "github.com/torfstack/kayvault/sql/gen"
)

type database struct {
	connStr string
	conn    *pgx.Conn
	tx      pgx.Tx
}

var _ Database = (*database)(nil)

func NewDatabase(connStr string) Database {
	return &database{connStr: connStr}
}

func (d *database) WithTx(ctx context.Context) (Database, Transaction) {
	conn, err := pgx.Connect(ctx, d.connStr)
	if err != nil {
		return nil, nil
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, nil
	}
	return &database{connStr: d.connStr, conn: conn, tx: tx}, &transaction{conn: conn, tx: tx}
}

func (d *database) CommitTransaction(ctx context.Context) error {
	if d.tx == nil {
		return nil
	}
	defer func(context.Context) {
		_ = (d.tx).Rollback(ctx)
		_ = (*d.conn).Close(ctx)
	}(ctx)
	return (d.tx).Commit(ctx)
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
	if secret.ID == nil {
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

func (d *database) InsertSecret(ctx context.Context, params sqlc.InsertSecretParams) error {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return err
	}
	return q.InsertSecret(ctx, params)
}

func (d *database) UpdateSecret(ctx context.Context, params sqlc.UpdateSecretParams) error {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return err
	}
	return q.UpdateSecret(ctx, params)
}

func startQuery(ctx context.Context, d *database) (*sqlc.Queries, error) {
	if d.tx != nil {
		return sqlc.New(d.tx), nil
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
