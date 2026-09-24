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
	return HashPasswordWithOptions(p, HashOptions{Salt: salt, Iterations: KeyDerivationIterations})
}

// HashPasswordWithOptions derives a hash using the caller-supplied salt and iterations.
// The caller is responsible for providing a cryptographically random salt.
func HashPasswordWithOptions(p Password, o HashOptions) (HashedPassword, error) {
	h, err := derivePasswordKey(p, o)
	if err != nil {
		return HashedPassword{}, err
	}
	return HashedPassword{
		Hash:           h,
		Salt:           o.Salt,
		IterationsUsed: o.Iterations,
	}, nil
}

func derivePasswordKey(password Password, options HashOptions) ([]byte, error) {
	if len(options.Salt) < KDFSaltLength || options.Iterations < KeyDerivationIterations ||
		options.Iterations > 2_000_000 {
		return nil, errors.New("invalid password derivation parameters")
	}
	return pbkdf2.Key(sha256.New, string(password), options.Salt, int(options.Iterations), AesKeyLengthInBytes)
}
