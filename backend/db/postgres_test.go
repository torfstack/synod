package db

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/torfstack/synod/backend/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	dbName     = "synod_test"
	dbUser     = "synod"
	dbPassword = "password"
)

var (
	pg *postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	var err error
	pg, err = setupTestContainer(context.Background())
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
	ctx := t.Context()
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
		createdUser, err := d.InsertUser(ctx, TestUser)
		assert.NoError(t, err)
		u, err := d.SelectUserByName(ctx, TestUser.Subject)
		assert.NoError(t, err)
		id := u.ID

		_, err = d.UpsertSecret(
			ctx, models.Secret{
				Value: "secret",
				Key:   "key",
				Url:   "url",
			}, createdUser.ID,
		)
		assert.NoError(t, err)

		secrets, err := d.SelectSecrets(ctx, id)
		assert.NoError(t, err)
		assert.Len(t, secrets, 1)
		assert.Equal(t, "secret", string(secrets[0].Value))
	}
}

func TestDatabase_UserHandling(t *testing.T) {
	ctx := t.Context()
	assert.NoError(t, pg.Restore(ctx))

	connStr, err := pg.ConnectionString(ctx)
	assert.NoError(t, err)
	d := NewDatabase(connStr)

	b, err := d.DoesUserExist(ctx, TestUser.Subject)
	assert.NoError(t, err)
	assert.False(t, b)

	createdUser, err := d.InsertUser(ctx, TestUser)
	assert.NoError(t, err)
	assert.Equal(t, TestUser.Subject, createdUser.Subject)
	assert.NotNil(t, createdUser.ID)

	b, err = d.DoesUserExist(ctx, TestUser.Subject)
	assert.NoError(t, err)
	assert.True(t, b)

	u, err := d.SelectUserByName(ctx, TestUser.Subject)
	assert.NoError(t, err)
	assert.Equal(t, TestUser.Subject, u.Subject)

	_, err = d.InsertUser(ctx, TestUser)
	assert.Error(t, err)
}

func TestDatabase_KeyHandling(t *testing.T) {
	ctx := t.Context()
	assert.NoError(t, pg.Restore(ctx))

	connStr, err := pg.ConnectionString(ctx)
	assert.NoError(t, err)
	d := NewDatabase(connStr)

	_, err = d.InsertKeys(ctx, models.UserKeyPair{})
	assert.Error(t, err)

	createdUser, err := d.InsertUser(ctx, TestUser)
	assert.NoError(t, err)

	createdKeys, err := d.InsertKeys(ctx, models.UserKeyPair{
		UserID: createdUser.ID,
		KeyPair: models.KeyPair{
			Public:  []byte{0x01, 0x02},
			Private: []byte{0x02, 0x03},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, createdKeys.UserID)
	assert.NotNil(t, createdKeys.ID)

	p, err := d.SelectPublicKey(ctx, createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x02}, p)
}

func TestDatabase_UserTransactionRollback(t *testing.T) {
	ctx := t.Context()
	assert.NoError(t, pg.Restore(ctx))

	connStr, err := pg.ConnectionString(ctx)
	assert.NoError(t, err)
	dd := NewDatabase(connStr)

	err = dd.WithTx(ctx, func(d Database) error {
		_, err = d.InsertUser(ctx, TestUser)
		assert.NoError(t, err)

		b, err := d.DoesUserExist(ctx, TestUser.Subject)
		assert.NoError(t, err)
		assert.True(t, b)
		return errors.New("trigger rollback")
	})
	require.Error(t, err)

	b, err := dd.DoesUserExist(ctx, TestUser.Subject)
	assert.NoError(t, err)
	assert.False(t, b)
}

func TestDatabase_UserTransactionCommit(t *testing.T) {
	ctx := t.Context()
	assert.NoError(t, pg.Restore(ctx))

	connStr, err := pg.ConnectionString(ctx)
	assert.NoError(t, err)
	dd := NewDatabase(connStr)

	err = dd.WithTx(ctx, func(d Database) error {
		_, err = d.InsertUser(ctx, TestUser)
		assert.NoError(t, err)

		b, err := d.DoesUserExist(ctx, TestUser.Subject)
		assert.NoError(t, err)
		assert.True(t, b)
		return nil
	})
	require.NoError(t, err)

	b, err := dd.DoesUserExist(ctx, TestUser.Subject)
	assert.NoError(t, err)
	assert.True(t, b)
}

func setupTestContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.WithSQLDriver("pgx"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}
	return postgresContainer, nil
}

var (
	TestUser = models.User{
		Subject:  "123-456-789-0",
		Email:    "tropfstein@gmail.com",
		FullName: "T.R. Opfstein",
	}
)
