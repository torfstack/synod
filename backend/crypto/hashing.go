package crypto

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

// Password is a mutable byte slice holding a plaintext password.
// Call Zero() to wipe the contents from memory when done.
type Password []byte

func (p Password) Zero() {
	for i := range p {
		p[i] = 0
	}
}

type HashedPassword struct {
	Hash           []byte
	Salt           []byte
	IterationsUsed int64
}

type HashOptions struct {
	Salt       []byte
	Iterations int64
}

func HashPassword(p Password) (HashedPassword, error) {
	salt := make([]byte, 32)
	read, err := rand.Read(salt)
	if err != nil {
		return HashedPassword{}, err
	}
	if read != 32 {
		return HashedPassword{}, errors.New("failed to generate 32 bytes of random salt")
	}
	h, _ := pbkdf2.Key(sha256.New, string(p), salt, KeyDerivationIterations, 32)
	return HashedPassword{
		Hash:           h,
		Salt:           salt,
		IterationsUsed: int64(KeyDerivationIterations),
	}, nil
}

// HashPasswordWithOptions derives a hash using the caller-supplied salt and iterations.
// The caller is responsible for providing a cryptographically random salt.
func HashPasswordWithOptions(p Password, o HashOptions) (HashedPassword, error) {
	h, _ := pbkdf2.Key(sha256.New, string(p), o.Salt, int(o.Iterations), 32)
	return HashedPassword{
		Hash:           h,
		Salt:           o.Salt,
		IterationsUsed: o.Iterations,
	}, nil
}
