package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type transaction struct {
	tx pgx.Tx
}

var _ Transaction = (*transaction)(nil)

func (t *transaction) Commit(ctx context.Context) {
	_ = t.tx.Commit(ctx)
}

func (t *transaction) Rollback(ctx context.Context) {
	_ = t.tx.Rollback(ctx)
}

func (t *transaction) SqlTx() pgx.Tx {
	return t.tx
}
