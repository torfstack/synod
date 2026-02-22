package crypto

import (
	"crypto/hpke"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AsymmetricCipher_Encrypt(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "can decrypt encrypted bytes with same cipher",
			run: func(t *testing.T) {
				a, err := NewAsymmetricCipher()
				require.NoError(t, err)

				b := []byte("synod testing encryptoin")
				e, err := a.Encrypt(b)
				require.NoError(t, err)

				d, err := a.Decrypt(e)
				require.NoError(t, err)
				assert.Equal(t, b, d)
			},
		},
		{
			name: "can decrypt encrypted bytes with recreated cipher",
			run: func(t *testing.T) {
				r, err := hpke.MLKEM768X25519().GenerateKey()
				require.NoError(t, err)

				a1, err := AsymmetricCipherFromPrivateKey(r)
				require.NoError(t, err)
				b := []byte("synod testing encryptoin")
				e, err := a1.Encrypt(b)
				require.NoError(t, err)

				a2, err := AsymmetricCipherFromPrivateKey(r)
				require.NoError(t, err)
				d, err := a2.Decrypt(e)
				require.NoError(t, err)
				assert.Equal(t, b, d)
			},
		},
		{
			name: "can serialize and deserialize asymmetric cipher",
			run: func(t *testing.T) {
				a, err := NewAsymmetricCipher()
				require.NoError(t, err)

				c, err := a.Encrypt([]byte("synod testing serialize"))
				require.NoError(t, err)

				b := a.Serialize()
				a2, err := AsymmetricCipherFromBytes(b)
				require.NoError(t, err)

				d, err := a2.Decrypt(c)
				require.NoError(t, err)
				assert.Equal(t, "synod testing serialize", string(d))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.run(t)
			},
		)
	}
}
