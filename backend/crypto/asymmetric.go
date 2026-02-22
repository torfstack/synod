package crypto

import (
	"crypto/hpke"
	"errors"
	"slices"
)

type AsymmetricCipher struct {
	publicKey  hpke.PublicKey
	privateKey hpke.PrivateKey
}

var (
	kem  = hpke.MLKEM768X25519()
	kdf  = hpke.HKDFSHA512()
	aead = hpke.AES256GCM()
)

func (a *AsymmetricCipher) Encrypt(plaintext []byte) ([]byte, error) {
	c, err := hpke.Seal(a.publicKey, kdf, aead, nil, plaintext)
	if err != nil {
		return nil, err
	}
	return slices.Concat(MarkerBytes, AsymmetricMarkerBytes, c), nil
}

func (a *AsymmetricCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < 9 {
		return nil, errors.New("ciphertext too short")
	}
	marker, algorithm := ciphertext[:4], ciphertext[4:8]
	if !slices.Equal(marker, MarkerBytes) || !slices.Equal(algorithm, AsymmetricMarkerBytes) {
		return nil, ErrCryptoInvalidMarker
	}
	p, err := hpke.Open(a.privateKey, kdf, aead, nil, ciphertext[8:])
	if err != nil {
		return nil, err
	}
	return p, nil
}

func NewAsymmetricCipher() (*AsymmetricCipher, error) {
	privateKey, err := hpke.MLKEM768X25519().GenerateKey()
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.PublicKey()
	return &AsymmetricCipher{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

func AsymmetricCipherFromPrivateKey(privateKey hpke.PrivateKey) (*AsymmetricCipher, error) {
	return &AsymmetricCipher{
		publicKey:  privateKey.PublicKey(),
		privateKey: privateKey,
	}, nil
}

func AsymmetricCipherFromBytes(b []byte) (*AsymmetricCipher, error) {
	priv, err := kem.NewPrivateKey(b)
	if err != nil {
		return nil, err
	}
	return AsymmetricCipherFromPrivateKey(priv)
}

func (a *AsymmetricCipher) Serialize() []byte {
	b, _ := a.privateKey.Bytes()
	return b
}
