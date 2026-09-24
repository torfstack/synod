package domain

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"

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
		selectSecretForManagerFn: func(context.Context, int64, int64) (models.AccessibleSecret, error) {
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

func TestShareSecretRejectsThresholdSecret(t *testing.T) {
	database := &mockDatabase{
		selectSecretForManagerFn: func(context.Context, int64, int64) (models.AccessibleSecret, error) {
			return models.AccessibleSecret{ID: 9, Threshold: 2}, nil
		},
		insertSecretAccessFn: func(context.Context, int64, int64, int64, []byte) error {
			t.Fatal("threshold secret received permanent access")
			return nil
		},
	}
	svc := &service{database: database, sessions: make(sessionStore)}

	err := svc.ShareSecret(context.Background(), 9, 1, "recipient", newTestCipher(t))

	require.ErrorIs(t, err, ErrThresholdSecretPermanentSharing)
}

func TestCreateThresholdSecretDistributesRecoverableShares(t *testing.T) {
	ciphers := map[int64]*crypto.AsymmetricCipher{}
	publicKeys := map[int64][]byte{}
	for _, id := range []int64{1, 2, 3} {
		cipher, err := crypto.NewAsymmetricCipher()
		require.NoError(t, err)
		ciphers[id] = cipher
		publicKeys[id], err = cipher.SerializePublicKey()
		require.NoError(t, err)
	}
	users := map[string]int64{"two": 2, "three": 3}
	shares := map[int64][]byte{}
	roles := map[int64]models.ThresholdRole{}
	var payload string
	database := &mockDatabase{
		selectUserBySharingIDFn: func(_ context.Context, sharingID string) (models.ExistingUser, error) {
			return models.ExistingUser{ID: users[sharingID]}, nil
		},
		selectPublicKeyFn: func(_ context.Context, userID int64) ([]byte, error) { return publicKeys[userID], nil },
		insertThresholdSecretFn: func(_ context.Context, secret models.EncryptedSecret, _ int64, threshold int) (int64, error) {
			require.Equal(t, 2, threshold)
			payload = secret.Value
			return 9, nil
		},
		insertThresholdShareFn: func(
			_ context.Context, secretID, userID int64, share []byte, role models.ThresholdRole,
		) error {
			require.Equal(t, int64(9), secretID)
			shares[userID] = share
			roles[userID] = role
			return nil
		},
	}
	svc := &service{database: database, sessions: make(sessionStore)}

	id, err := svc.CreateThresholdSecret(context.Background(), models.ThresholdSecretInput{
		Secret:    models.Secret{Key: "recovery", Value: "correct horse", Tags: []string{}},
		Threshold: 2, Participants: []models.ThresholdParticipantInput{
			{SharingID: "two", Role: models.ThresholdRoleHolder},
			{SharingID: "three", Role: models.ThresholdRoleHolder},
		},
	}, 1)
	require.NoError(t, err)
	require.Equal(t, int64(9), id)
	require.Equal(t, models.ThresholdRoleOwner, roles[1])
	require.Equal(t, models.ThresholdRoleHolder, roles[2])
	require.Equal(t, models.ThresholdRoleHolder, roles[3])

	decryptedShares := make([][]byte, 2)
	for i, userID := range []int64{1, 3} {
		decryptedShares[i], err = ciphers[userID].Decrypt(shares[userID])
		require.NoError(t, err)
	}
	key, err := crypto.CombineShares(decryptedShares)
	require.NoError(t, err)
	payloadCipher, err := crypto.SymmetricCipherFromKey(key)
	require.NoError(t, err)
	encoded, err := base64.StdEncoding.DecodeString(payload)
	require.NoError(t, err)
	plaintext, err := payloadCipher.Decrypt(encoded)
	require.NoError(t, err)
	var secret models.Secret
	require.NoError(t, json.Unmarshal(plaintext, &secret))
	require.Equal(t, "correct horse", secret.Value)
}

func TestSetThresholdParticipantsRejectsUserThatDisappearedBeforeSubmit(t *testing.T) {
	cipher, err := crypto.NewAsymmetricCipher()
	require.NoError(t, err)
	database := &mockDatabase{
		selectSecretForManagerFn: func(context.Context, int64, int64) (models.AccessibleSecret, error) {
			return models.AccessibleSecret{ID: 9, Threshold: 2, Role: models.ThresholdRoleOwner}, nil
		},
		selectThresholdParticipantsFn: func(context.Context, int64, int64) ([]models.ThresholdParticipant, error) {
			return []models.ThresholdParticipant{
				{ShareRecipient: models.ShareRecipient{Subject: "owner"}, Role: models.ThresholdRoleOwner},
				{ShareRecipient: models.ShareRecipient{Subject: "gone"}, Role: models.ThresholdRoleHolder},
			}, nil
		},
		selectThresholdParticipantKeysFn: func(context.Context, int64) ([]models.ParticipantKey, error) {
			return []models.ParticipantKey{{UserID: 1, Role: models.ThresholdRoleOwner}}, nil
		},
		selectUserBySharingIDFn: func(context.Context, string) (models.ExistingUser, error) {
			return models.ExistingUser{}, errors.New("user no longer exists")
		},
	}
	svc := &service{database: database, sessions: make(sessionStore)}

	err = svc.SetThresholdParticipants(context.Background(), 9, 1, models.ThresholdParticipantsInput{
		ExpectedParticipants: []models.ThresholdParticipantInput{{SharingID: "gone", Role: models.ThresholdRoleHolder}},
		Participants:         []models.ThresholdParticipantInput{{SharingID: "gone", Role: models.ThresholdRoleHolder}},
	}, cipher)

	require.ErrorIs(t, err, ErrThresholdParticipantsChanged)
}

func TestSetThresholdParticipantsRejectsStaleList(t *testing.T) {
	database := &mockDatabase{
		selectSecretForManagerFn: func(context.Context, int64, int64) (models.AccessibleSecret, error) {
			return models.AccessibleSecret{ID: 9, Threshold: 2, Role: models.ThresholdRoleOwner}, nil
		},
		selectThresholdParticipantsFn: func(context.Context, int64, int64) ([]models.ThresholdParticipant, error) {
			return []models.ThresholdParticipant{
				{ShareRecipient: models.ShareRecipient{Subject: "owner"}, Role: models.ThresholdRoleOwner},
				{ShareRecipient: models.ShareRecipient{Subject: "current"}, Role: models.ThresholdRoleHolder},
			}, nil
		},
		selectThresholdParticipantKeysFn: func(context.Context, int64) ([]models.ParticipantKey, error) {
			return []models.ParticipantKey{{UserID: 1, Role: models.ThresholdRoleOwner}}, nil
		},
	}
	svc := &service{database: database, sessions: make(sessionStore)}

	err := svc.SetThresholdParticipants(context.Background(), 9, 1, models.ThresholdParticipantsInput{
		ExpectedParticipants: []models.ThresholdParticipantInput{
			{SharingID: "removed", Role: models.ThresholdRoleHolder},
		},
		Participants: []models.ThresholdParticipantInput{
			{SharingID: "removed", Role: models.ThresholdRoleHolder},
		},
	}, newTestCipher(t))

	require.ErrorIs(t, err, ErrThresholdParticipantsChanged)
}

func TestGetUnlockResultGrantsEveryParticipantTenMinutesOfAccess(t *testing.T) {
	requester, err := crypto.NewAsymmetricCipher()
	require.NoError(t, err)
	dataKey, err := crypto.NewSymmetricKey()
	require.NoError(t, err)
	shares, err := crypto.SplitSecret(dataKey, 2, 3)
	require.NoError(t, err)
	encryptedContributions := make([][]byte, 2)
	for i := range encryptedContributions {
		encryptedContributions[i], err = requester.Encrypt(shares[i])
		require.NoError(t, err)
	}
	payloadCipher, err := crypto.SymmetricCipherFromKey(dataKey)
	require.NoError(t, err)
	plaintext, err := json.Marshal(models.Secret{Key: "shared", Value: "value", Tags: []string{}})
	require.NoError(t, err)
	payload, err := payloadCipher.Encrypt(plaintext)
	require.NoError(t, err)

	participantCiphers := map[int64]*crypto.AsymmetricCipher{}
	participantKeys := make([]models.ParticipantKey, 3)
	for i, userID := range []int64{1, 2, 3} {
		participantCiphers[userID], err = crypto.NewAsymmetricCipher()
		require.NoError(t, err)
		participantKeys[i].UserID = userID
		participantKeys[i].PublicKey, err = participantCiphers[userID].SerializePublicKey()
		require.NoError(t, err)
	}
	granted := map[int64][]byte{}
	completed := false
	now := time.Now()
	database := &mockDatabase{
		selectUnlockRequestFn: func(context.Context, int64, int64) (models.ThresholdSecret, error) {
			return models.ThresholdSecret{
				ID:               9,
				Threshold:        2,
				EncryptedPayload: []byte(base64.StdEncoding.EncodeToString(payload)),
			}, nil
		},
		selectUnlockContributionsFn: func(context.Context, int64) ([][]byte, error) {
			return encryptedContributions, nil
		},
		selectThresholdParticipantKeysFn: func(context.Context, int64) ([]models.ParticipantKey, error) {
			return participantKeys, nil
		},
		insertThresholdUnlockGrantFn: func(_ context.Context, requestID, userID int64, key []byte, expiresAt time.Time) error {
			require.Equal(t, int64(4), requestID)
			require.WithinDuration(t, now.Add(10*time.Minute), expiresAt, time.Second)
			granted[userID] = key
			return nil
		},
		completeUnlockRequestFn: func(context.Context, int64) error { completed = true; return nil },
	}
	svc := &service{database: database, sessions: make(sessionStore)}

	result, err := svc.GetUnlockResult(context.Background(), 4, 1, requester)
	require.NoError(t, err)
	require.True(t, result.Ready)
	require.Len(t, granted, 3)
	require.True(t, completed)
	for userID, encodedKey := range granted {
		wrappedKey, decodeErr := base64.StdEncoding.DecodeString(string(encodedKey))
		require.NoError(t, decodeErr)
		unwrappedKey, decryptErr := participantCiphers[userID].Decrypt(wrappedKey)
		require.NoError(t, decryptErr)
		require.Equal(t, dataKey, unwrappedKey)
	}
}

func TestStartUnlockRejectsSecondRequester(t *testing.T) {
	database := &mockDatabase{insertUnlockRequestFn: func(context.Context, int64, int64) (models.UnlockRequest, error) {
		return models.UnlockRequest{ID: 8, SecretID: 9, RequesterID: 1}, nil
	}}
	svc := &service{database: database, sessions: make(sessionStore)}

	_, err := svc.StartUnlock(context.Background(), 9, 2, newTestCipher(t))
	require.ErrorIs(t, err, ErrUnlockAlreadyActive)
}

func TestUpsertThresholdSecretUsesActiveGrantWithoutCreatingPermanentAccess(t *testing.T) {
	cipher := newTestCipher(t)
	key, err := crypto.NewSymmetricKey()
	require.NoError(t, err)
	wrappedKey, err := cipher.Encrypt(key)
	require.NoError(t, err)
	id := int64(9)
	var storedPayload string
	var storedSecret models.EncryptedSecret
	database := &mockDatabase{
		selectSecretForManagerFn: func(context.Context, int64, int64) (models.AccessibleSecret, error) {
			return models.AccessibleSecret{
				ID:               id,
				Envelope:         true,
				Threshold:        2,
				EncryptedDataKey: []byte(base64.StdEncoding.EncodeToString(wrappedKey)),
			}, nil
		},
		upsertSecretFn: func(_ context.Context, secret models.EncryptedSecret, _ int64) (models.EncryptedSecret, error) {
			storedPayload = secret.Value
			storedSecret = secret
			return secret, nil
		},
		insertSecretAccessFn: func(context.Context, int64, int64, int64, []byte) error {
			t.Fatal("threshold edit created permanent access")
			return nil
		},
	}
	svc := &service{database: database, sessions: make(sessionStore)}

	_, err = svc.UpsertSecret(
		context.Background(),
		models.Secret{
			ID:    &id,
			Key:   "updated",
			Value: "new value",
			Url:   "https://example.com",
			Tags:  []string{"recovery"},
		},
		1,
		cipher,
	)
	require.NoError(t, err)
	require.Equal(t, "updated", storedSecret.Key)
	require.Equal(t, "https://example.com", storedSecret.Url)
	require.Equal(t, []string{"recovery"}, storedSecret.Tags)
	encodedPayload, err := base64.StdEncoding.DecodeString(storedPayload)
	require.NoError(t, err)
	payloadCipher, err := crypto.SymmetricCipherFromKey(key)
	require.NoError(t, err)
	plaintext, err := payloadCipher.Decrypt(encodedPayload)
	require.NoError(t, err)
	var updated models.Secret
	require.NoError(t, json.Unmarshal(plaintext, &updated))
	require.Equal(t, "new value", updated.Value)
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
