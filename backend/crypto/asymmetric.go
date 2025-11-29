package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"slices"

	"github.com/torfstack/synod/backend/util"
)

type AsymmetricCipher struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (a *AsymmetricCipher) Encrypt(plaintext []byte) ([]byte, error) {
	s, err := NewSymmetricCipher()
	if err != nil {
		return nil, err
	}
	encryptedSymmetricKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, a.publicKey, s.key, nil)
	if err != nil {
		return nil, err
	}
	ciphertext, err := s.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	return slices.Concat(
		MarkerBytes, RsaOaepMarkerBytes,
		util.IntToBytes(uint32(len(encryptedSymmetricKey))), encryptedSymmetricKey,
		util.IntToBytes(uint32(len(ciphertext))), ciphertext,
	), nil
}

func (a *AsymmetricCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	b := bytes.NewBuffer(ciphertext)

	marker := b.Next(4)
	if !slices.Equal(marker, MarkerBytes) {
		return nil, ErrCryptoInvalidMarker
	}

	algorithm := b.Next(4)
	if !slices.Equal(algorithm, RsaOaepMarkerBytes) {
		return nil, ErrCryptoAlgorithmMarker
	}

	encryptedSymmetricKeyLen := util.BytesToInt(b.Next(4))
	encryptedSymmetricKey := b.Next(int(encryptedSymmetricKeyLen))

	symmetricKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, a.privateKey, encryptedSymmetricKey, nil)
	if err != nil {
		return nil, err
	}
	s, err := SymmetricCipherFromKey(symmetricKey)
	if err != nil {
		return nil, err
	}

	innerCiphertextLen := util.BytesToInt(b.Next(4))
	innerCiphertext := b.Next(int(innerCiphertextLen))

	plaintext, err := s.Decrypt(innerCiphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func NewAsymmetricCipher() (*AsymmetricCipher, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, RsaKeyLengthInBits)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.PublicKey
	return &AsymmetricCipher{
		publicKey:  &publicKey,
		privateKey: privateKey,
	}, nil
}

func AsymmetricCipherFromPublicKey(publicKey *rsa.PublicKey) (*AsymmetricCipher, error) {
	return &AsymmetricCipher{
		publicKey: publicKey,
	}, nil
}

func AsymmetricCipherFromPrivateKey(privateKey *rsa.PrivateKey) (*AsymmetricCipher, error) {
	return &AsymmetricCipher{
		publicKey:  privateKey.Public().(*rsa.PublicKey),
		privateKey: privateKey,
	}, nil
}

func AsymmetricCipherFromPrivateKeyBytes(b []byte) (*AsymmetricCipher, error) {
	priv, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	priv.Precompute()
	return AsymmetricCipherFromPrivateKey(priv)
}

func (a *AsymmetricCipher) Serialize() []byte {
	return x509.MarshalPKCS1PrivateKey(a.privateKey)
}
