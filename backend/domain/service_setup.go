package domain

import (
	"context"
	"crypto/subtle"
	"errors"

	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/models"
)

var (
	KeyDerivationSalt = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

var _ SetupService = &service{}

func (s *service) IsUserSetup(ctx context.Context, session Session) (bool, error) {
	return s.database.HasKeys(ctx, session.UserID)
}

func (s *service) SetupUserPlain(ctx context.Context, session Session) error {
	a, err := crypto.NewAsymmetricCipher()
	if err != nil {
		return err
	}
	_, err = s.database.InsertKeys(ctx, models.UserKeyPair{
		UserID:      session.UserID,
		Type:        models.KeyTypeRsa,
		KeyMaterial: a.Serialize(),
	})
	if err != nil {
		return err
	}
	session.Cipher = a
	s.sessions[session.SessionID] = session
	return err
}

func (s *service) SetupUserWithPassword(ctx context.Context, session Session, password string) error {
	a, err := crypto.NewAsymmetricCipher()
	if err != nil {
		return err
	}
	p, err := crypto.SymmetricCipherFromPassword([]byte(password))
	if err != nil {
		return err
	}
	encrypted, err := p.Encrypt(a.Serialize())
	if err != nil {
		return err
	}

	hashedPassword, err := crypto.HashPassword([]byte(password))
	if err != nil {
		return err
	}
	dbPassword, err := s.database.InsertPassword(ctx, models.HashedPassword{
		Hash:       hashedPassword.Hash,
		Salt:       hashedPassword.Salt,
		Iterations: hashedPassword.IterationsUsed,
	})
	if err != nil {
		return err
	}

	_, err = s.database.InsertKeys(ctx, models.UserKeyPair{
		UserID:      session.UserID,
		PasswordID:  dbPassword.ID,
		Type:        models.KeyTypeRsa,
		KeyMaterial: encrypted,
	})
	if err != nil {
		return err
	}

	session.Cipher = a
	s.sessions[session.SessionID] = session

	return err
}

func (s *service) UnsealWithPassword(ctx context.Context, session *Session, password string) error {
	if session.Cipher != nil {
		return nil
	}

	key, err := s.database.SelectKeys(ctx, session.UserID)
	if err != nil {
		return err
	}
	if key.PasswordID == nil {
		return errors.New("no password associated with keys in database")
	}

	dbPassword, err := s.database.SelectPassword(ctx, *key.PasswordID)
	if err != nil {
		return err
	}

	hashedPassword, err := crypto.HashPasswordWithOptions([]byte(password), crypto.HashOptions{
		Salt:       dbPassword.Salt,
		Iterations: dbPassword.Iterations,
	})
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare(dbPassword.Hash, hashedPassword.Hash) == 0 {
		return errors.New("password hash mismatch")
	}

	p, err := crypto.SymmetricCipherFromPassword([]byte(password))
	if err != nil {
		return err
	}

	decryptedPrivateKey, err := p.Decrypt(key.KeyMaterial)
	if err != nil {
		return err
	}

	a, err := crypto.AsymmetricCipherFromPrivateKeyBytes(decryptedPrivateKey)
	if err != nil {
		return err
	}

	session.Cipher = a
	s.sessions[session.SessionID] = *session

	return nil
}
