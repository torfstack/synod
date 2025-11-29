package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"slices"

	"github.com/torfstack/synod/backend/util"
)

type SymmetricCipher struct {
	cipher cipher.AEAD
	key    []byte
}

func (s *SymmetricCipher) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, s.cipher.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	sealed := s.cipher.Seal(nil, nonce, plaintext, nil)
	return slices.Concat(
		MarkerBytes,
		AesGcmMarkerBytes,
		util.IntToBytes(uint32(len(nonce))),
		nonce,
		util.IntToBytes(uint32(len(sealed))),
		sealed,
	), nil
}

func (s *SymmetricCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	b := bytes.NewBuffer(ciphertext)

	marker := b.Next(4)
	if !slices.Equal(marker, MarkerBytes) {
		return nil, ErrCryptoInvalidMarker
	}

	algorithm := b.Next(4)
	if !slices.Equal(algorithm, AesGcmMarkerBytes) {
		return nil, ErrCryptoAlgorithmMarker
	}

	nonceLen := util.BytesToInt(b.Next(4))
	nonce := b.Next(int(nonceLen))
	sealedLen := util.BytesToInt(b.Next(4))
	sealed := b.Next(int(sealedLen))

	return s.cipher.Open(nil, nonce, sealed, nil)
}

func NewSymmetricKey() ([]byte, error) {
	key := make([]byte, AesKeyLengthInBytes)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func SymmetricCipherFromKey(key []byte) (*SymmetricCipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	c, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &SymmetricCipher{c, key}, nil
}

func SymmetricCipherFromPassword(password []byte) (*SymmetricCipher, error) {
	derivedKey, _ := pbkdf2.Key(sha256.New, string(password), KeyDerivationSalt, 600000, 32)
	return SymmetricCipherFromKey(derivedKey)
}

func NewSymmetricCipher() (*SymmetricCipher, error) {
	key, err := NewSymmetricKey()
	if err != nil {
		return nil, err
	}
	return SymmetricCipherFromKey(key)
}
