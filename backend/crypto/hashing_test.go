package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPasswordWithOptionsRejectsInvalidParameters(t *testing.T) {
	valid := HashOptions{Salt: []byte("1234567890abcdef"), Iterations: KeyDerivationIterations}
	for _, change := range []func(*HashOptions){
		func(o *HashOptions) { o.Salt = []byte("short") },
		func(o *HashOptions) { o.Iterations = 0 },
		func(o *HashOptions) { o.Iterations = 1 << 30 },
	} {
		invalid := valid
		change(&invalid)
		_, err := HashPasswordWithOptions(Password("password"), invalid)
		require.Error(t, err)
	}
}

func Test_HashPassword(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "hashing twice yields different results based on different salt",
			run: func(t *testing.T) {
				b := Password("synod testing password")
				h1, err := HashPassword(b)
				assert.NoError(t, err)
				assert.NotEqual(t, b, h1.Hash)

				h2, err := HashPassword(b)
				assert.NoError(t, err)
				assert.NotEqual(t, h1.Hash, h2.Hash)
				assert.NotEqual(t, h1.Salt, h2.Salt)
				assert.Equal(t, h1.IterationsUsed, h2.IterationsUsed)
			},
		},
		{
			name: "hashing twice with the same options yields the same result",
			run: func(t *testing.T) {
				b := Password("synod password hashing")
				salt := []byte("1234567890abcdef")
				h1, err := HashPasswordWithOptions(b, HashOptions{
					Salt:       salt,
					Iterations: KeyDerivationIterations,
				})
				assert.NoError(t, err)
				assert.NotEqual(t, b, h1.Hash)
				assert.Equal(t, salt, h1.Salt)
				assert.Equal(t, int64(KeyDerivationIterations), h1.IterationsUsed)

				h2, err := HashPasswordWithOptions(b, HashOptions{
					Salt:       salt,
					Iterations: KeyDerivationIterations,
				})
				assert.NoError(t, err)
				assert.Equal(t, h1.Hash, h2.Hash)

				h3, err := HashPasswordWithOptions(b, HashOptions{
					Salt:       []byte("1234567890abcdeg"),
					Iterations: KeyDerivationIterations,
				})
				assert.NoError(t, err)
				assert.NotEqual(t, h1.Hash, h3.Hash)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}
