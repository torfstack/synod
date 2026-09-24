package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"slices"

	"github.com/torfstack/synod/backend/util"
)

type SymmetricCipher struct {
	cipher cipher.AEAD
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
	if len(ciphertext) < 8 {
		return nil, ErrCryptoInvalidCiphertext
	}
	if !slices.Equal(ciphertext[:4], MarkerBytes) {
		return nil, ErrCryptoInvalidMarker
	}
	if !slices.Equal(ciphertext[4:8], AesGcmMarkerBytes) {
		return nil, ErrCryptoAlgorithmMarker
	}
	if len(ciphertext) < 12 {
		return nil, ErrCryptoInvalidCiphertext
	}
	nonceLen := util.BytesToInt(ciphertext[8:12])
	if nonceLen != uint32(s.cipher.NonceSize()) || len(ciphertext) < 12+int(nonceLen)+4 {
		return nil, ErrCryptoInvalidCiphertext
	}
	nonceEnd := 12 + int(nonceLen)
	sealedLen := util.BytesToInt(ciphertext[nonceEnd : nonceEnd+4])
	sealed := ciphertext[nonceEnd+4:]
	if sealedLen < uint32(s.cipher.Overhead()) || uint64(sealedLen) != uint64(len(sealed)) {
		return nil, ErrCryptoInvalidCiphertext
	}
	return s.cipher.Open(nil, ciphertext[12:nonceEnd], sealed, nil)
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
	return &SymmetricCipher{cipher: c}, nil
}

func SymmetricCipherFromPasswordWithSalt(password Password, salt []byte) (*SymmetricCipher, error) {
	derivedKey, err := derivePasswordKey(password, HashOptions{Salt: salt, Iterations: KeyDerivationIterations})
	if err != nil {
		return nil, err
	}
	defer clear(derivedKey)
	return SymmetricCipherFromKey(derivedKey)
}

func NewSymmetricCipher() (*SymmetricCipher, error) {
	key, err := NewSymmetricKey()
	if err != nil {
		return nil, err
	}
	return SymmetricCipherFromKey(key)
}
