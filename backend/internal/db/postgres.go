package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"main/internal/config"
	sqlc "main/sql/gen"
)

type postgres struct {
	cfg config.DBConfig
}

func NewDatabase(cfg config.DBConfig) Database {
	return &postgres{
		cfg: cfg,
	}
}

func (p *postgres) Connect(ctx context.Context) (Connection, error) {
	conn, err := pgx.Connect(ctx, p.cfg.ConnectionString())
	if err != nil {
		return nil, err
	}
	return connection{conn: conn}, nil
}

type connection struct {
	conn *pgx.Conn
}

func (c connection) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

func (c connection) Queries() Queries {
	return sqlc.New(c.conn)
}
