package db

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"log"
	"os"
	"testing"
	"time"

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
	require.NoError(t, err)
	d, err := NewDatabase(ctx, connStr)
	require.NoError(t, err)
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
			ctx, models.EncryptedSecret{
				Value: "encrypted_secret",
				Key:   "key",
				Url:   "url",
			}, createdUser.ID,
		)
		assert.NoError(t, err)

		secrets, err := d.SelectSecrets(ctx, id)
		assert.NoError(t, err)
		assert.Len(t, secrets, 1)
		assert.Equal(t, "encrypted_secret", string(secrets[0].Value))
	}
}

func TestDatabase_DeleteExpiredUnlockRequestsCascadesRelatedData(t *testing.T) {
	ctx := t.Context()
	require.NoError(t, pg.Restore(ctx))
	connStr, err := pg.ConnectionString(ctx)
	require.NoError(t, err)
	databaseInterface, err := NewDatabase(ctx, connStr)
	require.NoError(t, err)
	database := databaseInterface.(*database)
	user, err := database.InsertUser(ctx, TestUser)
	require.NoError(t, err)
	secretID, err := database.InsertThresholdSecret(
		ctx,
		models.EncryptedSecret{Value: "payload", Key: "key"},
		user.ID,
		2,
	)
	require.NoError(t, err)
	request, err := database.InsertUnlockRequest(ctx, secretID, user.ID)
	require.NoError(t, err)
	_, err = database.pool.Exec(
		ctx,
		"INSERT INTO unlock_contributions (request_id, user_id, encrypted_share) VALUES ($1, $2, $3)",
		request.ID,
		user.ID,
		[]byte("share"),
	)
	require.NoError(t, err)
	_, err = database.pool.Exec(
		ctx,
		"INSERT INTO threshold_unlock_grants (request_id, user_id, encrypted_data_key, expires_at) VALUES ($1, $2, $3, NOW() - INTERVAL '1 minute')",
		request.ID,
		user.ID,
		[]byte("key"),
	)
	require.NoError(t, err)
	require.NoError(t, database.CompleteUnlockRequest(ctx, request.ID))

	deleted, err := database.DeleteExpiredUnlockRequests(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(1), deleted)
	var requests, contributions, grants int
	require.NoError(t, database.pool.QueryRow(ctx, "SELECT COUNT(*) FROM unlock_requests").Scan(&requests))
	require.NoError(t, database.pool.QueryRow(ctx, "SELECT COUNT(*) FROM unlock_contributions").Scan(&contributions))
	require.NoError(t, database.pool.QueryRow(ctx, "SELECT COUNT(*) FROM threshold_unlock_grants").Scan(&grants))
	require.Zero(t, requests)
	require.Zero(t, contributions)
	require.Zero(t, grants)
}

func TestDatabase_AllowsOnlyOneActiveUnlockRequestPerSecret(t *testing.T) {
	ctx := t.Context()
	require.NoError(t, pg.Restore(ctx))
	connStr, err := pg.ConnectionString(ctx)
	require.NoError(t, err)
	databaseInterface, err := NewDatabase(ctx, connStr)
	require.NoError(t, err)
	database := databaseInterface.(*database)
	owner, err := database.InsertUser(ctx, TestUser)
	require.NoError(t, err)
	participant, err := database.InsertUser(
		ctx,
		models.User{Subject: "participant", Email: "participant@example.com", FullName: "Participant"},
	)
	require.NoError(t, err)
	secretID, err := database.InsertThresholdSecret(
		ctx,
		models.EncryptedSecret{Value: "payload", Key: "key"},
		owner.ID,
		2,
	)
	require.NoError(t, err)
	require.NoError(t, database.InsertThresholdShare(
		ctx, secretID, participant.ID, []byte("share"), models.ThresholdRoleHolder,
	))

	first, err := database.InsertUnlockRequest(ctx, secretID, owner.ID)
	require.NoError(t, err)
	second, err := database.InsertUnlockRequest(ctx, secretID, participant.ID)
	require.NoError(t, err)
	require.Equal(t, first.ID, second.ID)
	require.Equal(t, owner.ID, second.RequesterID)

	_, err = database.pool.Exec(
		ctx,
		"UPDATE unlock_requests SET expires_at = NOW() - INTERVAL '1 minute' WHERE id = $1",
		first.ID,
	)
	require.NoError(t, err)
	replacement, err := database.InsertUnlockRequest(ctx, secretID, participant.ID)
	require.NoError(t, err)
	require.NotEqual(t, first.ID, replacement.ID)
	require.Equal(t, participant.ID, replacement.RequesterID)
}

func TestDatabase_ThresholdOwnerCanUseActiveUnlockGrant(t *testing.T) {
	ctx := t.Context()
	require.NoError(t, pg.Restore(ctx))
	connStr, err := pg.ConnectionString(ctx)
	require.NoError(t, err)
	databaseInterface, err := NewDatabase(ctx, connStr)
	require.NoError(t, err)
	database := databaseInterface.(*database)
	owner, err := database.InsertUser(ctx, TestUser)
	require.NoError(t, err)
	secretID, err := database.InsertThresholdSecret(
		ctx,
		models.EncryptedSecret{Value: "payload", Key: "key"},
		owner.ID,
		2,
	)
	require.NoError(t, err)
	request, err := database.InsertUnlockRequest(ctx, secretID, owner.ID)
	require.NoError(t, err)
	encryptedKey := []byte("encrypted key")
	require.NoError(
		t,
		database.InsertThresholdUnlockGrant(ctx, request.ID, owner.ID, encryptedKey, time.Now().Add(time.Minute)),
	)

	secret, err := database.SelectSecretForManager(ctx, secretID, owner.ID)
	require.NoError(t, err)
	require.Equal(t, encryptedKey, secret.EncryptedDataKey)
	require.Equal(t, 2, secret.Threshold)
	require.NoError(t, database.CompleteUnlockRequest(ctx, request.ID))
	_, err = database.SelectUnlockRequest(ctx, request.ID, owner.ID)
	require.Error(t, err)
}

func TestDatabase_UserHandling(t *testing.T) {
	ctx := t.Context()
	assert.NoError(t, pg.Restore(ctx))

	connStr, err := pg.ConnectionString(ctx)
	assert.NoError(t, err)
	d, err := NewDatabase(ctx, connStr)
	require.NoError(t, err)

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
	d, err := NewDatabase(ctx, connStr)
	require.NoError(t, err)

	createdUser, err := d.InsertUser(ctx, TestUser)
	assert.NoError(t, err)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	createdKeys, err := d.InsertKeys(
		ctx, models.UserKeyPair{
			UserID:      createdUser.ID,
			Type:        models.KeyTypeRsa,
			KeyMaterial: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, createdKeys.UserID)
	assert.NotNil(t, createdKeys.ID)

	keyPair, err := d.SelectKeys(ctx, createdUser.ID)
	assert.NoError(t, err)

	privParsed, err := x509.ParsePKCS1PrivateKey(keyPair.KeyMaterial)
	require.NoError(t, err)
	assert.True(t, priv.Equal(privParsed), "stored RSA key must match the original")

	assert.Equal(t, createdUser.ID, keyPair.UserID)
}

func TestDatabase_UserTransactionRollback(t *testing.T) {
	ctx := t.Context()
	assert.NoError(t, pg.Restore(ctx))

	connStr, err := pg.ConnectionString(ctx)
	assert.NoError(t, err)
	dd, err := NewDatabase(ctx, connStr)
	require.NoError(t, err)

	err = dd.WithTx(
		ctx, func(d Database) error {
			_, err = d.InsertUser(ctx, TestUser)
			assert.NoError(t, err)

			b, err := d.DoesUserExist(ctx, TestUser.Subject)
			assert.NoError(t, err)
			assert.True(t, b)
			return errors.New("trigger rollback")
		},
	)
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
	dd, err := NewDatabase(ctx, connStr)
	require.NoError(t, err)

	err = dd.WithTx(
		ctx, func(d Database) error {
			_, err = d.InsertUser(ctx, TestUser)
			assert.NoError(t, err)

			b, err := d.DoesUserExist(ctx, TestUser.Subject)
			assert.NoError(t, err)
			assert.True(t, b)
			return nil
		},
	)
	require.NoError(t, err)

	b, err := dd.DoesUserExist(ctx, TestUser.Subject)
	assert.NoError(t, err)
	assert.True(t, b)
}

func setupTestContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:18-alpine",
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
