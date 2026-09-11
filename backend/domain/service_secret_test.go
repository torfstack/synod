package domain

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/models"
)

func TestShareSecretWrapsDataKeyForRecipient(t *testing.T) {
	owner, err := crypto.NewAsymmetricCipher()
	require.NoError(t, err)
	recipient, err := crypto.NewAsymmetricCipher()
	require.NoError(t, err)
	recipientPublicKey, err := recipient.SerializePublicKey()
	require.NoError(t, err)
	dataKey, err := crypto.NewSymmetricKey()
	require.NoError(t, err)
	wrappedForOwner, err := owner.Encrypt(dataKey)
	require.NoError(t, err)

	var wrappedForRecipient []byte
	db := &mockDatabase{
		selectSecretForOwnerFn: func(context.Context, int64, int64) (models.AccessibleSecret, error) {
			return models.AccessibleSecret{
				ID:               9,
				OwnerID:          1,
				Envelope:         true,
				EncryptedDataKey: []byte(base64.StdEncoding.EncodeToString(wrappedForOwner)),
			}, nil
		},
		selectUserBySharingIDFn: func(context.Context, string) (models.ExistingUser, error) {
			return models.ExistingUser{ID: 2}, nil
		},
		selectPublicKeyFn: func(context.Context, int64) ([]byte, error) { return recipientPublicKey, nil },
		insertSecretAccessFn: func(_ context.Context, _, userID, grantedBy int64, key []byte) error {
			require.Equal(t, int64(2), userID)
			require.Equal(t, int64(1), grantedBy)
			wrappedForRecipient = key
			return nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	require.NoError(t, svc.ShareSecret(context.Background(), 9, 1, "recipient", owner))
	wrapped, err := base64.StdEncoding.DecodeString(string(wrappedForRecipient))
	require.NoError(t, err)
	unwrapped, err := recipient.Decrypt(wrapped)
	require.NoError(t, err)
	require.Equal(t, dataKey, unwrapped)
}

func TestGetSecretsDecryptsSharedEnvelope(t *testing.T) {
	recipient, err := crypto.NewAsymmetricCipher()
	require.NoError(t, err)
	dataKey, err := crypto.NewSymmetricKey()
	require.NoError(t, err)
	payloadCipher, err := crypto.SymmetricCipherFromKey(dataKey)
	require.NoError(t, err)
	payload, err := json.Marshal(models.Secret{Key: "database", Value: "password", Tags: []string{}})
	require.NoError(t, err)
	encryptedPayload, err := payloadCipher.Encrypt(payload)
	require.NoError(t, err)
	wrappedKey, err := recipient.Encrypt(dataKey)
	require.NoError(t, err)
	db := &mockDatabase{selectAccessibleSecretsFn: func(context.Context, int64) ([]models.AccessibleSecret, error) {
		return []models.AccessibleSecret{{
			ID: 4, OwnerID: 1, Owned: false, Envelope: true,
			EncryptedPayload: []byte(base64.StdEncoding.EncodeToString(encryptedPayload)),
			EncryptedDataKey: []byte(base64.StdEncoding.EncodeToString(wrappedKey)),
		}}, nil
	}}
	svc := &service{database: db, sessions: make(sessionStore)}

	secrets, err := svc.GetSecrets(context.Background(), 2, recipient)
	require.NoError(t, err)
	require.Len(t, secrets, 1)
	require.Equal(t, "password", secrets[0].Value)
	require.False(t, secrets[0].Owned)
}

// newTestCipher creates a real AsymmetricCipher for use in tests.
func newTestCipher(t *testing.T) *crypto.AsymmetricCipher {
	t.Helper()
	c, err := crypto.NewAsymmetricCipher()
	require.NoError(t, err)
	return c
}

// ---- GetSecrets ----

func TestGetSecrets_NilCipher_ReturnsError(t *testing.T) {
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}

	secrets, err := svc.GetSecrets(context.Background(), 1, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cipher")
	assert.Empty(t, secrets)
}

func TestGetSecrets_DatabaseError_Propagates(t *testing.T) {
	dbErr := errors.New("db down")
	db := &mockDatabase{
		selectSecretsFn: func(_ context.Context, _ int64) ([]models.EncryptedSecret, error) {
			return nil, dbErr
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}
	cipher := newTestCipher(t)

	_, err := svc.GetSecrets(context.Background(), 1, cipher)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
}

func TestGetSecrets_DecryptsAndReturnsSecrets(t *testing.T) {
	cipher := newTestCipher(t)

	// Encrypt a plaintext so we can store a valid ciphertext in the mock
	plaintext := "super-secret-value"
	cipherBytes, err := cipher.Encrypt([]byte(plaintext))
	require.NoError(t, err)
	encoded := base64.StdEncoding.EncodeToString(cipherBytes)

	db := &mockDatabase{
		selectSecretsFn: func(_ context.Context, _ int64) ([]models.EncryptedSecret, error) {
			return []models.EncryptedSecret{
				{Value: encoded, Key: "k1", Url: "u1", Tags: []string{"tag1"}},
			}, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}

	secrets, err := svc.GetSecrets(context.Background(), 1, cipher)

	require.NoError(t, err)
	require.Len(t, secrets, 1)
	assert.Equal(t, plaintext, secrets[0].Value)
	assert.Equal(t, "k1", secrets[0].Key)
}

func TestGetSecrets_EmptyList_ReturnsEmpty(t *testing.T) {
	db := &mockDatabase{
		selectSecretsFn: func(_ context.Context, _ int64) ([]models.EncryptedSecret, error) {
			return []models.EncryptedSecret{}, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}
	cipher := newTestCipher(t)

	secrets, err := svc.GetSecrets(context.Background(), 1, cipher)

	require.NoError(t, err)
	assert.Empty(t, secrets)
}

// ---- UpsertSecret ----

func TestUpsertSecret_NilCipher_ReturnsError(t *testing.T) {
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}

	_, err := svc.UpsertSecret(context.Background(), models.Secret{}, 1, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cipher")
}

func TestUpsertSecret_DatabaseError_ReturnsError(t *testing.T) {
	db := &mockDatabase{
		upsertSecretFn: func(_ context.Context, _ models.EncryptedSecret, _ int64) (models.EncryptedSecret, error) {
			return models.EncryptedSecret{}, errors.New("db write failed")
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}
	cipher := newTestCipher(t)

	_, err := svc.UpsertSecret(context.Background(), models.Secret{Value: "v", Key: "k"}, 1, cipher)

	require.Error(t, err)
}

func TestUpsertSecret_Success_ReturnsEncryptedSecret(t *testing.T) {
	id := int64(99)
	db := &mockDatabase{
		upsertSecretFn: func(_ context.Context, s models.EncryptedSecret, _ int64) (models.EncryptedSecret, error) {
			s.ID = &id
			return s, nil
		},
	}
	svc := &service{database: db, sessions: make(sessionStore)}
	cipher := newTestCipher(t)

	result, err := svc.UpsertSecret(
		context.Background(),
		models.Secret{Value: "my-secret", Key: "k", Url: "u"},
		42,
		cipher,
	)

	require.NoError(t, err)
	require.NotNil(t, result.ID)
	assert.Equal(t, id, *result.ID)
	// The stored value should be base64-encoded ciphertext, not the plaintext
	assert.NotEqual(t, "my-secret", result.Value)
}
