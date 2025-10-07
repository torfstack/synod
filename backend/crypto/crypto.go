package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"slices"

	"github.com/torfstack/synod/backend/models"
	"github.com/torfstack/synod/backend/util"
)

var (
	MarkerBytes = []byte{0x64, 0x61, 0x74, 0x71}

	RsaOaepMarkerBytes = []byte{0x00, 0x00, 0x00, 0x01}
	AesGcmMarkerBytes  = []byte{0x00, 0x00, 0x00, 0x02}

	RsaKeyLengthInBits  = 2048
	AesKeyLengthInBytes = 32

	ErrCryptoInvalidMarker   = errors.New("invalid encryption marker")
	ErrCryptoAlgorithmMarker = errors.New("invalid algorithm marker")
)

type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

type symmetricCipher struct {
	cipher cipher.AEAD
}

func (s *symmetricCipher) Encrypt(plaintext []byte) ([]byte, error) {
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

func (s *symmetricCipher) Decrypt(ciphertext []byte) ([]byte, error) {
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

func NewPublicKeyPair() (models.KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, RsaKeyLengthInBits)
	if err != nil {
		return models.KeyPair{}, err
	}
	publicKey := privateKey.PublicKey
	return models.KeyPair{
		Public:  publicKey,
		Private: *privateKey,
	}, nil
}

func NewSymmetricKey() ([]byte, error) {
	key := make([]byte, AesKeyLengthInBytes)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func SymmetricCipherFromKey(key []byte) (Cipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	c, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &symmetricCipher{c}, nil
}

func NewSymmetricCipher() (Cipher, error) {
	key, err := NewSymmetricKey()
	if err != nil {
		return nil, err
	}
	return SymmetricCipherFromKey(key)
}
