package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HashPassword(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "hashing twice yields different results based on different salt",
			run: func(t *testing.T) {
				b := []byte("synod testing password")
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
				b := []byte("synod password hashing")
				h1, err := HashPasswordWithOptions(b, HashOptions{
					Salt:       []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
					Iterations: 500000,
				})
				assert.NoError(t, err)
				assert.NotEqual(t, b, h1.Hash)
				assert.Equal(t, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, h1.Salt)
				assert.Equal(t, int64(500000), h1.IterationsUsed)

				h2, err := HashPasswordWithOptions(b, HashOptions{
					Salt:       []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
					Iterations: 500000,
				})
				assert.NoError(t, err)
				assert.Equal(t, h1.Hash, h2.Hash)

				h3, err := HashPasswordWithOptions(b, HashOptions{
					Salt:       []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x08},
					Iterations: 500000,
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
