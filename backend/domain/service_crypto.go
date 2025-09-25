package domain

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/torfstack/synod/backend/models"
)

var _ CryptoService = &service{}

func (s *service) GenerateKeyPair() (models.KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return models.KeyPair{}, err
	}
	return models.KeyPair{
		Public:  pub,
		Private: priv,
	}, nil
}
