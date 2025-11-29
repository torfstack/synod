package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SymmetricCipher_Encrypt(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "encrypted bytes can be decrypted on the same cipher",
			run: func(t *testing.T) {
				cipher, err := NewSymmetricCipher()
				require.NoError(t, err)

				encrypted, err := cipher.Encrypt([]byte("synod"))
				require.NoError(t, err)

				decrypted, err := cipher.Decrypt(encrypted)
				require.NoError(t, err)
				require.Equal(t, "synod", string(decrypted))
			},
		},
		{
			name: "encrypted bytes can be decrypted with a different cipher on the same key",
			run: func(t *testing.T) {
				key, err := NewSymmetricKey()
				require.NoError(t, err)

				cipher1, err := SymmetricCipherFromKey(key)
				require.NoError(t, err)
				encrypted, err := cipher1.Encrypt([]byte("synod"))
				require.NoError(t, err)

				cipher2, err := SymmetricCipherFromKey(key)
				require.NoError(t, err)

				decrypted, err := cipher2.Decrypt(encrypted)
				require.NoError(t, err)
				require.Equal(t, "synod", string(decrypted))
			},
		},
		{
			name: "encrypted bytes can not be decrypted with a different cipher on a different key",
			run: func(t *testing.T) {
				key1, err := NewSymmetricKey()
				require.NoError(t, err)
				cipher1, err := SymmetricCipherFromKey(key1)
				require.NoError(t, err)

				encrypted, err := cipher1.Encrypt([]byte("synod"))
				require.NoError(t, err)

				key2, err := NewSymmetricKey()
				require.NoError(t, err)
				require.NotEqual(t, key1, key2)
				cipher2, err := SymmetricCipherFromKey(key2)
				require.NoError(t, err)

				_, err = cipher2.Decrypt(encrypted)
				require.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}
