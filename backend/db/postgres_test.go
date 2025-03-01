package db

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	sqlc "github.com/torfstack/kayvault/sql/gen"
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
	connStr, err := pg.ConnectionString(context.Background())
	if err != nil {
		log.Printf("failed to get connection string: %s", err)
		return
	}
	err = Migrate(context.Background(), connStr, WithMigrationsDir("../../sql/migrations"))
	if err != nil {
		log.Printf("failed to migrate database: %s", err)
		return
	}
	err = pg.Snapshot(context.Background())
	if err != nil {
		log.Printf("failed to snapshot container: %s", err)
		return
	}
	os.Exit(m.Run())
}

func TestDatabase_SecretHandling(t *testing.T) {
	ctx := context.Background()
	assert.NoError(t, pg.Restore(ctx))

	connStr, err := pg.ConnectionString(ctx)
	assert.NoError(t, err)
	d := NewDatabase(connStr)
	{
		secrets, err := d.SelectSecrets(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, secrets, 0)
	}
	{
		assert.NoError(t, d.InsertUser(ctx, TestUserParams))
		u, err := d.SelectUserByName(ctx, TestUserParams.Subject)
		assert.NoError(t, err)
		id := u.ID

		err = d.InsertSecret(
			ctx, sqlc.InsertSecretParams{
				Value:  []byte("secret"),
				Key:    "key",
				Url:    "url",
				UserID: id,
			},
		)
		assert.NoError(t, err)

		secrets, err := d.SelectSecrets(ctx, id)
		assert.NoError(t, err)
		assert.Len(t, secrets, 1)
		assert.Equal(t, "secret", string(secrets[0].Value))
	}
}

func TestDatabase_UserHandling(t *testing.T) {
	ctx := context.Background()
	assert.NoError(t, pg.Restore(ctx))

	connStr, err := pg.ConnectionString(ctx)
	assert.NoError(t, err)
	d := NewDatabase(connStr)

	b, err := d.DoesUserExist(ctx, TestUserParams.Subject)
	assert.NoError(t, err)
	assert.False(t, b)

	err = d.InsertUser(ctx, TestUserParams)
	assert.NoError(t, err)

	b, err = d.DoesUserExist(ctx, TestUserParams.Subject)
	assert.NoError(t, err)
	assert.True(t, b)

	u, err := d.SelectUserByName(ctx, TestUserParams.Subject)
	assert.NoError(t, err)
	assert.Equal(t, TestUserParams.Subject, u.Subject)

	err = d.InsertUser(ctx, TestUserParams)
	assert.Error(t, err)
}

func TestDatabase_UserTransactionRollback(t *testing.T) {
	ctx := context.Background()
	assert.NoError(t, pg.Restore(ctx))

	connStr, err := pg.ConnectionString(ctx)
	assert.NoError(t, err)
	dd := NewDatabase(connStr)

	d, tx := dd.WithTx(ctx)
	{
		err := d.InsertUser(ctx, TestUserParams)
		assert.NoError(t, err)

		b, err := d.DoesUserExist(ctx, TestUserParams.Subject)
		assert.NoError(t, err)
		assert.True(t, b)
	}
	tx.Rollback(ctx)

	b, err := dd.DoesUserExist(ctx, TestUserParams.Subject)
	assert.NoError(t, err)
	assert.False(t, b)
}

func TestDatabase_UserTransactionCommit(t *testing.T) {
	ctx := context.Background()
	assert.NoError(t, pg.Restore(ctx))

	connStr, err := pg.ConnectionString(ctx)
	assert.NoError(t, err)
	dd := NewDatabase(connStr)

	d, tx := dd.WithTx(ctx)
	{
		err := d.InsertUser(ctx, TestUserParams)
		assert.NoError(t, err)

		b, err := d.DoesUserExist(ctx, TestUserParams.Subject)
		assert.NoError(t, err)
		assert.True(t, b)
	}
	tx.Commit(ctx)

	b, err := dd.DoesUserExist(ctx, TestUserParams.Subject)
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
		postgres.WithSQLDriver("pgx"),
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

var (
	TestUserParams = sqlc.InsertUserParams{
		Subject:  "123-456-789-0",
		Email:    "tropfstein@gmail.com",
		FullName: "T.R. Opfstein",
	}
)
