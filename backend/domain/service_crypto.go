package domain

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"slices"

	"github.com/torfstack/synod/backend/models"
)

var _ CryptoService = &service{}

var (
	MarkerBytes        = []byte{0x64, 0x61, 0x74, 0x71}
	RsaOaepMarkerBytes = []byte{0x00, 0x00, 0x00, 0x01}
	AesGcmMarkerBytes  = []byte{0x00, 0x00, 0x00, 0x02}
)

func (s *service) GenerateKeyPair() (models.KeyPair, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return models.KeyPair{}, err
	}
	pub := priv.PublicKey
	return models.KeyPair{
		Public:  pub,
		Private: *priv,
	}, nil
}

func (s *service) EncryptSecret(
	_ context.Context,
	secret models.Secret,
	key *rsa.PublicKey,
) (models.EncryptedSecret, error) {
	aesKey, err := newAesKey()
	if err != nil {
		return models.EncryptedSecret{}, err
	}
	gcmCipher, err := aesCipher(aesKey)
	if err != nil {
		return models.EncryptedSecret{}, err
	}
	encryptedAesKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, aesKey, nil)
	if err != nil {
		return models.EncryptedSecret{}, err
	}
	l := intToBytes(len(encryptedAesKey))
	ciphertext, nonce, err := encrypt([]byte(secret.Value), gcmCipher)
	if err != nil {
		return models.EncryptedSecret{}, err
	}
	nonceL := intToBytes(len(nonce))
	cipherTextL := intToBytes(len(ciphertext))
	encrypted := slices.Concat(MarkerBytes, RsaOaepMarkerBytes, l, encryptedAesKey, nonceL, nonce, cipherTextL, ciphertext)
	return models.EncryptedSecret{
		ID:    secret.ID,
		Value: base64.StdEncoding.EncodeToString(encrypted),
		Key:   secret.Key,
		Url:   secret.Url,
		Tags:  secret.Tags,
	}, nil
}

func (s *service) DecryptSecret(
	_ context.Context,
	secret models.EncryptedSecret,
	key *rsa.PrivateKey,
) (models.Secret, error) {
	b, err := base64.StdEncoding.DecodeString(secret.Value)
	if err != nil {
		return models.Secret{}, err
	}
	header := b[:8]
	if !slices.Equal(header, slices.Concat(MarkerBytes, RsaOaepMarkerBytes)) {
		return models.Secret{}, errors.New("invalid encryption header bytes")
	}
	encryptedAesKeyLength := bytesToInt(b[8:12])
	encryptedAesKey := b[12 : 12+encryptedAesKeyLength]
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, encryptedAesKey, nil)
	if err != nil {
		return models.Secret{}, err
	}
	nonceLength := bytesToInt(b[12+encryptedAesKeyLength : 12+encryptedAesKeyLength+4])
	nonce := b[12+encryptedAesKeyLength+4 : 12+encryptedAesKeyLength+4+nonceLength]
	cipherTextLength := bytesToInt(b[12+encryptedAesKeyLength+4+nonceLength : 12+encryptedAesKeyLength+4+nonceLength+4])
	cipherText := b[12+encryptedAesKeyLength+4+nonceLength+4 : 12+encryptedAesKeyLength+4+nonceLength+4+cipherTextLength]

	gcmCipher, err := aesCipher(aesKey)
	if err != nil {
		return models.Secret{}, err
	}
	decrypted, err := decrypt(cipherText, gcmCipher, nonce)
	if err != nil {
		return models.Secret{}, err
	}

	return models.Secret{
		ID:    secret.ID,
		Value: string(decrypted),
		Key:   secret.Key,
		Url:   secret.Url,
		Tags:  secret.Tags,
	}, nil
}

func newAesKey() ([]byte, error) {
	keyBytes := make([]byte, 32)
	read, err := rand.Reader.Read(keyBytes)
	if err != nil {
		return nil, err
	}
	if read != len(keyBytes) {
		return nil, errors.New("failed to read random bytes")
	}
	return keyBytes, nil
}

func aesCipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func encrypt(plaintext []byte, cipher cipher.AEAD) (encrypted []byte, nonce []byte, err error) {
	nonce = make([]byte, cipher.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return
	}
	return cipher.Seal(nil, nonce, plaintext, nil), nonce, nil
}

func decrypt(ciphertext []byte, cipher cipher.AEAD, nonce []byte) ([]byte, error) {
	return cipher.Open(nil, nonce, ciphertext, nil)
}

func intToBytes(length int) []byte {
	buffer := new(bytes.Buffer)
	// err is always nil
	_ = binary.Write(buffer, binary.BigEndian, uint32(length))
	b := buffer.Bytes()
	return []byte{b[0], b[1], b[2], b[3]}
}

func bytesToInt(b []byte) int {
	return int(binary.BigEndian.Uint32(b))
}
