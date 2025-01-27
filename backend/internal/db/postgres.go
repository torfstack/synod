package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/torfstack/kayvault/internal/config"
	sqlc "github.com/torfstack/kayvault/sql/gen"
)

type Database interface {
	WithTx(ctx context.Context) (Database, Transaction)

	DoesUserExist(ctx context.Context, username string) (bool, error)
	InsertUser(ctx context.Context, username string) error
	SelectUserByName(ctx context.Context, username string) (sqlc.User, error)
	InsertSecret(ctx context.Context, params sqlc.InsertSecretParams) error
	UpdateSecret(ctx context.Context, params sqlc.UpdateSecretParams) error
	SelectSecrets(ctx context.Context, userID int32) ([]sqlc.Secret, error)
}

type Transaction interface {
	Commit(ctx context.Context)

	// Rollback is a no-op if the transaction has already been committed
	Rollback(ctx context.Context)
}

type database struct {
	cfg  config.DBConfig
	conn *pgx.Conn
	tx   pgx.Tx
}

type transaction struct {
	conn *pgx.Conn
	tx   pgx.Tx
}

func NewDatabase(cfg config.DBConfig) Database {
	return &database{cfg: cfg}
}

func (d *database) WithTx(ctx context.Context) (Database, Transaction) {
	conn, err := pgx.Connect(ctx, d.cfg.ConnectionString())
	if err != nil {
		return nil, nil
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, nil
	}
	return &database{cfg: d.cfg, conn: conn, tx: tx}, &transaction{conn: conn, tx: tx}
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

func (d *database) InsertUser(ctx context.Context, username string) error {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return err
	}
	return q.InsertUser(ctx, username)
}

func (d *database) SelectUserByName(ctx context.Context, username string) (sqlc.User, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return sqlc.User{}, err
	}
	return q.SelectUserByName(ctx, username)
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

func (d *database) SelectSecrets(ctx context.Context, userID int32) ([]sqlc.Secret, error) {
	q, err := startQuery(ctx, d)
	defer endQuery(ctx, d)
	if err != nil {
		return []sqlc.Secret{}, err
	}
	return q.SelectSecrets(ctx, userID)
}

func startQuery(ctx context.Context, d *database) (*sqlc.Queries, error) {
	if d.tx != nil {
		return sqlc.New(d.tx), nil
	}
	if d.conn == nil {
		conn, err := pgx.Connect(ctx, d.cfg.ConnectionString())
		d.conn = conn
		return nil, err
	}
	return sqlc.New(d.conn), nil
}

func endQuery(ctx context.Context, d *database) {
	if d.conn != nil && d.tx == nil {
		_ = (*d.conn).Close(ctx)
	}
}

func (t *transaction) Commit(ctx context.Context) {
	_ = t.tx.Commit(ctx)
	_ = t.tx.Conn().Close(ctx)
}

func (t *transaction) Rollback(ctx context.Context) {
	_ = t.tx.Rollback(ctx)
	_ = t.tx.Conn().Close(ctx)
}
