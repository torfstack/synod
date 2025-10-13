package domain

import (
	"context"
	"encoding/base64"

	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/models"
)

var _ CryptoService = &service{}

var (
	MarkerBytes        = []byte{0x64, 0x61, 0x74, 0x71}
	RsaOaepMarkerBytes = []byte{0x00, 0x00, 0x00, 0x01}
)

func (s *service) EncryptSecret(
	_ context.Context,
	secret models.Secret,
	c *crypto.AsymmetricCipher,
) (models.EncryptedSecret, error) {
	ciphertext, err := c.Encrypt([]byte(secret.Value))
	if err != nil {
		return models.EncryptedSecret{}, err
	}
	return models.EncryptedSecret{
		ID:    secret.ID,
		Value: base64.StdEncoding.EncodeToString(ciphertext),
		Key:   secret.Key,
		Url:   secret.Url,
		Tags:  secret.Tags,
	}, nil
}

func (s *service) DecryptSecret(
	_ context.Context,
	secret models.EncryptedSecret,
	c *crypto.AsymmetricCipher,
) (models.Secret, error) {
	b, err := base64.StdEncoding.DecodeString(secret.Value)
	if err != nil {
		return models.Secret{}, err
	}
	decrypted, err := c.Decrypt(b)
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
