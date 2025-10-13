package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"slices"

	"github.com/torfstack/synod/backend/util"
	"golang.org/x/crypto/pbkdf2"
)

var (
	MarkerBytes = []byte{0x64, 0x61, 0x74, 0x71}

	RsaOaepMarkerBytes = []byte{0x00, 0x00, 0x00, 0x01}
	AesGcmMarkerBytes  = []byte{0x00, 0x00, 0x00, 0x02}

	RsaKeyLengthInBits  = 2048
	AesKeyLengthInBytes = 32

	KeyDerivationSalt = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	Iterations        = 600000

	ErrCryptoInvalidMarker   = errors.New("invalid encryption marker")
	ErrCryptoAlgorithmMarker = errors.New("invalid algorithm marker")
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
	derivedKey := pbkdf2.Key(password, KeyDerivationSalt, 600000, 32, sha256.New)
	return SymmetricCipherFromKey(derivedKey)
}

func NewSymmetricCipher() (*SymmetricCipher, error) {
	key, err := NewSymmetricKey()
	if err != nil {
		return nil, err
	}
	return SymmetricCipherFromKey(key)
}

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

type HashedPassword struct {
	Hash           []byte
	Salt           []byte
	IterationsUsed int64
}

type HashOptions struct {
	Salt       []byte
	Iterations int64
}

func HashPassword(p []byte) (HashedPassword, error) {
	salt := make([]byte, 32)
	read, err := rand.Read(salt)
	if err != nil {
		return HashedPassword{}, err
	}
	if read != 32 {
		return HashedPassword{}, errors.New("failed to generate 32 bytes of random salt")
	}
	h := pbkdf2.Key(p, salt, Iterations, 32, sha256.New)
	return HashedPassword{
		Hash:           h,
		Salt:           salt,
		IterationsUsed: int64(Iterations),
	}, nil
}

func HashPasswordWithOptions(p []byte, o HashOptions) HashedPassword {
	salt := make([]byte, 32)
	read, err := rand.Read(salt)
	if err != nil {
		return HashedPassword{}
	}
	if read != 32 {
		return HashedPassword{}
	}
	h := pbkdf2.Key(p, o.Salt, int(o.Iterations), 32, sha256.New)
	return HashedPassword{
		Hash:           h,
		Salt:           o.Salt,
		IterationsUsed: o.Iterations,
	}
}
