package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type transaction struct {
	conn *pgx.Conn
	tx   pgx.Tx
}

var _ Transaction = (*transaction)(nil)

func (t *transaction) Commit(ctx context.Context) {
	_ = t.tx.Commit(ctx)
	_ = t.tx.Conn().Close(ctx)
}

func (t *transaction) Rollback(ctx context.Context) {
	_ = t.tx.Rollback(ctx)
	_ = t.tx.Conn().Close(ctx)
}
