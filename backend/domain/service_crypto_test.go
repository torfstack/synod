package domain

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/torfstack/synod/backend/models"
)

// ---- EncryptSecret ----

func TestEncryptSecret_ProducesBase64EncodedValue(t *testing.T) {
	cipher := newTestCipher(t)
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}

	id := int64(1)
	secret := models.Secret{
		ID:    &id,
		Value: "my-plaintext",
		Key:   "k",
		Url:   "u",
		Tags:  []string{"t"},
	}

	encrypted, err := svc.EncryptSecret(context.Background(), secret, cipher)

	require.NoError(t, err)
	// Value should be valid base64
	decoded, err := base64.StdEncoding.DecodeString(encrypted.Value)
	require.NoError(t, err)
	assert.NotEmpty(t, decoded)
	// Value should NOT be the plaintext
	assert.NotEqual(t, "my-plaintext", encrypted.Value)
	// Other fields should be preserved
	assert.Equal(t, secret.ID, encrypted.ID)
	assert.Equal(t, secret.Key, encrypted.Key)
	assert.Equal(t, secret.Url, encrypted.Url)
	assert.Equal(t, secret.Tags, encrypted.Tags)
}

// ---- DecryptSecret ----

func TestDecryptSecret_RecoverPlaintext(t *testing.T) {
	cipher := newTestCipher(t)
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}

	plain := models.Secret{Value: "round-trip-value", Key: "k", Url: "u", Tags: []string{"a", "b"}}

	encrypted, err := svc.EncryptSecret(context.Background(), plain, cipher)
	require.NoError(t, err)

	decrypted, err := svc.DecryptSecret(context.Background(), encrypted, cipher)
	require.NoError(t, err)

	assert.Equal(t, plain.Value, decrypted.Value)
	assert.Equal(t, plain.Key, decrypted.Key)
	assert.Equal(t, plain.Url, decrypted.Url)
	assert.Equal(t, plain.Tags, decrypted.Tags)
}

func TestDecryptSecret_InvalidBase64_ReturnsError(t *testing.T) {
	cipher := newTestCipher(t)
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}

	encrypted := models.EncryptedSecret{Value: "!!!not-valid-base64!!!"}
	_, err := svc.DecryptSecret(context.Background(), encrypted, cipher)

	require.Error(t, err)
}

func TestDecryptSecret_ValidBase64ButWrongCipher_ReturnsError(t *testing.T) {
	cipher1 := newTestCipher(t)
	cipher2 := newTestCipher(t)
	svc := &service{database: &mockDatabase{}, sessions: make(sessionStore)}

	plain := models.Secret{Value: "secret"}
	encrypted, err := svc.EncryptSecret(context.Background(), plain, cipher1)
	require.NoError(t, err)

	_, err = svc.DecryptSecret(context.Background(), encrypted, cipher2)
	require.Error(t, err)
}
