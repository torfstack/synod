package domain

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"slices"

	"github.com/torfstack/synod/backend/models"
)

var _ CryptoService = &service{}

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

func (s *service) EncryptSecret(_ context.Context, secret models.Secret, key rsa.PublicKey) (models.EncryptedSecret, error) {
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key, []byte(secret.Value), nil)
	if err != nil {
		return models.EncryptedSecret{}, err
	}
	encrypted = append(append(MARKER_BYTES, RSA_MARKER_BYTES...), []byte(encrypted)...)
	return models.EncryptedSecret{
		ID:    secret.ID,
		Value: base64.StdEncoding.EncodeToString(encrypted),
		Key:   secret.Key,
		Url:   secret.Url,
		Tags:  secret.Tags,
	}, nil
}

func (s *service) DecryptSecret(_ context.Context, secret models.EncryptedSecret, key rsa.PrivateKey) (models.Secret, error) {
	bytes, err := base64.StdEncoding.DecodeString(secret.Value)
	if err != nil {
		return models.Secret{}, err
	}
	header, encryptedBytes := bytes[:8], bytes[8:]
	if !slices.Equal(header, append(MARKER_BYTES, RSA_MARKER_BYTES...)) {
		return models.Secret{}, errors.New("invalid encryption header bytes")
	}
	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &key, encryptedBytes, nil)
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

func (s *service) GetPublicKey(ctx context.Context, userID int64) (rsa.PublicKey, error) {
	return s.database.SelectPublicKey(ctx, userID)
}

func (s *service) GetPrivateKey(ctx context.Context, userID int64) (rsa.PrivateKey, error) {
	return s.database.SelectPrivateKey(ctx, userID)
}

var MARKER_BYTES = []byte{0x64, 0x61, 0x74, 0x71}
var RSA_MARKER_BYTES = []byte{0x00, 0x00, 0x00, 0x01}
