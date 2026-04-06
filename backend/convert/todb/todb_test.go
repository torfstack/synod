package todb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/torfstack/synod/backend/models"
)

// ---- Secret ----

func TestSecret_MapsAllFields(t *testing.T) {
	id := int64(5)
	in := models.Secret{
		ID:    &id,
		Value: "plaintext",
		Key:   "my-key",
		Url:   "https://example.com",
		Tags:  []string{"x", "y"},
	}
	out := Secret(in)

	assert.Equal(t, id, out.ID)
	assert.Equal(t, []byte("plaintext"), out.Value)
	assert.Equal(t, "my-key", out.Key)
	assert.Equal(t, "https://example.com", out.Url)
	assert.Equal(t, "x,y", out.Tags)
}

func TestSecret_NoTags(t *testing.T) {
	id := int64(1)
	in := models.Secret{ID: &id, Tags: []string{}}
	out := Secret(in)
	assert.Equal(t, "", out.Tags)
}

// ---- InsertSecretParams ----

func TestInsertSecretParams_MapsAllFields(t *testing.T) {
	in := models.EncryptedSecret{
		Value: "enc-value",
		Key:   "k",
		Url:   "u",
		Tags:  []string{"t1", "t2"},
	}
	out := InsertSecretParams(in, 99)

	assert.Equal(t, []byte("enc-value"), out.Value)
	assert.Equal(t, "k", out.Key)
	assert.Equal(t, "u", out.Url)
	assert.Equal(t, "t1,t2", out.Tags)
	assert.Equal(t, int64(99), out.UserID)
}

func TestInsertSecretParams_EmptyTags(t *testing.T) {
	in := models.EncryptedSecret{Tags: []string{}}
	out := InsertSecretParams(in, 1)
	assert.Equal(t, "", out.Tags)
}

// ---- UpdateSecretParams ----

func TestUpdateSecretParams_MapsAllFields(t *testing.T) {
	id := int64(7)
	in := models.EncryptedSecret{
		ID:    &id,
		Value: "updated",
		Key:   "k2",
		Url:   "u2",
		Tags:  []string{"a"},
	}
	out := UpdateSecretParams(in, 11)

	assert.Equal(t, id, out.ID)
	assert.Equal(t, []byte("updated"), out.Value)
	assert.Equal(t, "k2", out.Key)
	assert.Equal(t, "u2", out.Url)
	assert.Equal(t, "a", out.Tags)
	assert.Equal(t, int64(11), out.UserID)
}

// ---- InsertUserParams ----

func TestInsertUserParams_MapsAllFields(t *testing.T) {
	in := models.User{
		Subject:  "sub-abc",
		Email:    "user@test.com",
		FullName: "Bob Jones",
	}
	out := InsertUserParams(in)

	assert.Equal(t, "sub-abc", out.Subject)
	assert.Equal(t, "user@test.com", out.Email)
	assert.Equal(t, "Bob Jones", out.FullName)
}

// ---- InsertKeysParams ----

func TestInsertKeysParams_WithoutPassword(t *testing.T) {
	in := models.UserKeyPair{
		UserID:      3,
		Type:        models.KeyTypeRsa,
		KeyMaterial: []byte{0xAB, 0xCD},
		PasswordID:  nil,
	}
	out := InsertKeysParams(in)

	assert.Equal(t, int64(3), out.UserID)
	assert.Equal(t, int32(models.KeyTypeRsa), out.Type)
	assert.Equal(t, []byte{0xAB, 0xCD}, out.KeyMaterial)
	assert.False(t, out.PasswordID.Valid)
}

func TestInsertKeysParams_WithPassword(t *testing.T) {
	pwID := int64(55)
	in := models.UserKeyPair{
		UserID:     1,
		Type:       models.KeyTypeRsa,
		PasswordID: &pwID,
	}
	out := InsertKeysParams(in)

	require.True(t, out.PasswordID.Valid)
	assert.Equal(t, pwID, out.PasswordID.Int64)
}

// ---- InsertPasswordParams ----

func TestInsertPasswordParams_MapsAllFields(t *testing.T) {
	in := models.HashedPassword{
		Hash:       []byte("h"),
		Salt:       []byte("s"),
		Iterations: 200000,
	}
	out := InsertPasswordParams(in)

	assert.Equal(t, []byte("h"), out.Hash)
	assert.Equal(t, []byte("s"), out.Salt)
	assert.Equal(t, int64(200000), out.Iterations)
}

// ---- tagsString round-trip ----

func TestTagsString_SingleTag(t *testing.T) {
	in := models.Secret{ID: new(int64), Tags: []string{"only"}}
	out := Secret(in)
	assert.Equal(t, "only", out.Tags)
}

func TestTagsString_MultipleTagsNoTrailingComma(t *testing.T) {
	in := models.EncryptedSecret{Tags: []string{"a", "b", "c"}}
	out := InsertSecretParams(in, 1)
	assert.Equal(t, "a,b,c", out.Tags)
	// Ensure no trailing comma
	assert.NotEqual(t, "a,b,c,", out.Tags)
}
