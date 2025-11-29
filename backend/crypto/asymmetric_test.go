package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert.NoError(t, err)

				b := []byte("synod testing encryptoin")
				e, err := a.Encrypt(b)
				assert.NoError(t, err)
				assert.NotEqual(t, b, e)

				d, err := a.Decrypt(e)
				assert.NoError(t, err)
				assert.Equal(t, b, d)
			},
		},
		{
			name: "can decrypt encrypted bytes with recreated cipher",
			run: func(t *testing.T) {
				r, err := rsa.GenerateKey(rand.Reader, RsaKeyLengthInBits)
				assert.NoError(t, err)

				a1, err := AsymmetricCipherFromPrivateKey(r)
				assert.NoError(t, err)
				b := []byte("synod testing encryptoin")
				e, err := a1.Encrypt(b)
				assert.NoError(t, err)

				a2, err := AsymmetricCipherFromPrivateKey(r)
				assert.NoError(t, err)
				d, err := a2.Decrypt(e)
				assert.NoError(t, err)
				assert.Equal(t, b, d)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}
