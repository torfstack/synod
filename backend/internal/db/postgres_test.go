package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

const (
	dbName     = "kayvault_test"
	dbUser     = "kayvault"
	dbPassword = "password"
)

var (
	pg *postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	var err error
	pg, err = setupTestContainer()
	if err != nil {
		log.Printf("failed to setup test container")
		return
	}
	defer func() {
		if err = testcontainers.TerminateContainer(pg); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	connStr, _ := pg.ConnectionString(context.Background())
	_ = Migrate(context.Background(), connStr, WithMigrationsDir("../../sql/migrations"))
	err = pg.Snapshot(context.Background())
	if err != nil {
		log.Printf("failed to snapshot container: %s", err)
		return
	}
	m.Run()
}

func TestDatabase_UserHandling(t *testing.T) {
	ctx := context.Background()
	_ = pg.Restore(ctx)

	connStr, _ := pg.ConnectionString(ctx)
	d := NewDatabase(connStr)

	b, err := d.DoesUserExist(ctx, "tropfstein")
	assert.NoError(t, err)
	assert.False(t, b)

	err = d.InsertUser(ctx, "tropfstein")
	assert.NoError(t, err)

	b, err = d.DoesUserExist(ctx, "tropfstein")
	assert.NoError(t, err)
	assert.True(t, b)

	u, err := d.SelectUserByName(ctx, "tropfstein")
	assert.NoError(t, err)
	assert.Equal(t, "tropfstein", u.Username)

	err = d.InsertUser(ctx, "tropfstein")
	assert.Error(t, err)
}

func TestDatabase_UserTransactionRollback(t *testing.T) {
	ctx := context.Background()
	_ = pg.Restore(ctx)

	connStr, _ := pg.ConnectionString(ctx)
	dd := NewDatabase(connStr)

	d, tx := dd.WithTx(ctx)
	{
		err := d.InsertUser(ctx, "tropfstein")
		assert.NoError(t, err)

		b, err := d.DoesUserExist(ctx, "tropfstein")
		assert.NoError(t, err)
		assert.True(t, b)
	}
	tx.Rollback(ctx)

	b, err := dd.DoesUserExist(ctx, "tropfstein")
	assert.NoError(t, err)
	assert.False(t, b)
}

func TestDatabase_UserTransactionCommit(t *testing.T) {
	ctx := context.Background()
	_ = pg.Restore(ctx)

	connStr, _ := pg.ConnectionString(ctx)
	dd := NewDatabase(connStr)

	d, tx := dd.WithTx(ctx)
	{
		err := d.InsertUser(ctx, "tropfstein")
		assert.NoError(t, err)

		b, err := d.DoesUserExist(ctx, "tropfstein")
		assert.NoError(t, err)
		assert.True(t, b)
	}
	tx.Commit(ctx)

	b, err := dd.DoesUserExist(ctx, "tropfstein")
	assert.NoError(t, err)
	assert.True(t, b)
}

func setupTestContainer() (*postgres.PostgresContainer, error) {
	ctx := context.Background()
	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}
	return postgresContainer, nil
}
