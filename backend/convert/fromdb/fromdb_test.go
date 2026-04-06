package fromdb

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/torfstack/synod/backend/models"
	sqlc "github.com/torfstack/synod/sql/gen"
)

// ---- Secret / Secrets ----

func TestSecret_MapsAllFields(t *testing.T) {
	id := int64(7)
	in := sqlc.Secret{
		ID:    id,
		Value: []byte("cipher-text"),
		Key:   "my-key",
		Url:   "https://example.com",
		Tags:  "a,b,c",
	}
	out := Secret(in)

	require.NotNil(t, out.ID)
	assert.Equal(t, id, *out.ID)
	assert.Equal(t, "cipher-text", out.Value)
	assert.Equal(t, "my-key", out.Key)
	assert.Equal(t, "https://example.com", out.Url)
	assert.Equal(t, []string{"a", "b", "c"}, out.Tags)
}

func TestSecret_EmptyTags(t *testing.T) {
	in := sqlc.Secret{Tags: ""}
	out := Secret(in)
	assert.Equal(t, []string{}, out.Tags)
}

func TestSecret_SingleTag(t *testing.T) {
	in := sqlc.Secret{Tags: "solo"}
	out := Secret(in)
	assert.Equal(t, []string{"solo"}, out.Tags)
}

func TestSecrets_PreservesOrder(t *testing.T) {
	in := []sqlc.Secret{
		{ID: 1, Value: []byte("v1"), Key: "k1"},
		{ID: 2, Value: []byte("v2"), Key: "k2"},
		{ID: 3, Value: []byte("v3"), Key: "k3"},
	}
	out := Secrets(in)
	require.Len(t, out, 3)
	assert.Equal(t, "v1", out[0].Value)
	assert.Equal(t, "v2", out[1].Value)
	assert.Equal(t, "v3", out[2].Value)
}

func TestSecrets_Empty(t *testing.T) {
	out := Secrets([]sqlc.Secret{})
	assert.Empty(t, out)
}

// ---- User ----

func TestUser_MapsAllFields(t *testing.T) {
	in := sqlc.User{
		ID:       42,
		Subject:  "sub-123",
		Email:    "user@example.com",
		FullName: "Alice Smith",
	}
	out := User(in)

	assert.Equal(t, int64(42), out.ID)
	assert.Equal(t, "sub-123", out.Subject)
	assert.Equal(t, "user@example.com", out.Email)
	assert.Equal(t, "Alice Smith", out.FullName)
}

// ---- KeyPair ----

func TestKeyPair_WithoutPassword(t *testing.T) {
	id := int64(10)
	in := sqlc.Key{
		ID:          id,
		UserID:      5,
		Type:        int32(models.KeyTypeRsa),
		KeyMaterial: []byte{0x01, 0x02},
		PasswordID:  pgtype.Int8{Valid: false},
	}
	out := KeyPair(in)

	require.NotNil(t, out.ID)
	assert.Equal(t, id, *out.ID)
	assert.Equal(t, int64(5), out.UserID)
	assert.Equal(t, models.KeyTypeRsa, out.Type)
	assert.Equal(t, []byte{0x01, 0x02}, out.KeyMaterial)
	assert.Nil(t, out.PasswordID)
}

func TestKeyPair_WithPassword(t *testing.T) {
	passwordID := int64(99)
	in := sqlc.Key{
		ID:          1,
		UserID:      2,
		Type:        int32(models.KeyTypeRsa),
		KeyMaterial: []byte{0xAA},
		PasswordID:  pgtype.Int8{Int64: passwordID, Valid: true},
	}
	out := KeyPair(in)

	require.NotNil(t, out.PasswordID)
	assert.Equal(t, passwordID, *out.PasswordID)
}

// ---- HashedPassword ----

func TestHashedPassword_MapsAllFields(t *testing.T) {
	id := int64(3)
	in := sqlc.Password{
		ID:         id,
		Hash:       []byte("hash-bytes"),
		Salt:       []byte("salt-bytes"),
		Iterations: 100000,
	}
	out := HashedPassword(in)

	require.NotNil(t, out.ID)
	assert.Equal(t, id, *out.ID)
	assert.Equal(t, []byte("hash-bytes"), out.Hash)
	assert.Equal(t, []byte("salt-bytes"), out.Salt)
	assert.Equal(t, int64(100000), out.Iterations)
}
