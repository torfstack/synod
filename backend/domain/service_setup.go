package domain

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/x509"
	"errors"
	"slices"

	"github.com/torfstack/synod/backend/models"
	"golang.org/x/crypto/pbkdf2"
)

var (
	KeyDerivationSalt = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

var _ SetupService = &service{}

func (s *service) IsUserSetup(ctx context.Context, session Session) (bool, error) {
	return s.database.HasKeys(ctx, session.UserID)
}

func (s *service) SetupUserPlain(ctx context.Context, session Session) error {
	keyPair, err := s.GenerateKeyPair()
	if err != nil {
		return err
	}
	_, err = s.database.InsertKeys(ctx, models.UserKeyPair{
		UserID:  session.UserID,
		Type:    models.KeyTypeRsa,
		Public:  x509.MarshalPKCS1PublicKey(&keyPair.Public),
		Private: x509.MarshalPKCS1PrivateKey(&keyPair.Private),
	})
	session.PrivateKey = &keyPair.Private
	s.sessions[session.SessionID] = session
	return err
}

func (s *service) SetupUserWithPassword(ctx context.Context, session Session, password string) error {
	keyPair, err := s.GenerateKeyPair()
	if err != nil {
		return err
	}
	derivedKey := pbkdf2.Key([]byte(password), KeyDerivationSalt, 600000, 32, sha256.New)
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(&keyPair.Private)
	gcm, err := aesCipher(derivedKey)
	if err != nil {
		return err
	}
	encrypted, nonce, err := encrypt(privateKeyBytes, gcm)
	if err != nil {
		return err
	}
	nonceL := intToBytes(len(nonce))
	encryptedL := intToBytes(len(encrypted))
	ciphertext := slices.Concat(MarkerBytes, AesGcmMarkerBytes, nonceL, nonce, encryptedL, encrypted)

	salt := make([]byte, 8)
	read, err := rand.Reader.Read(salt)
	if err != nil {
		return err
	}
	if read != 8 {
		return errors.New("failed to generate random salt")
	}
	hashedPassword := pbkdf2.Key([]byte(password), salt, 600000, 32, sha256.New)
	dbPassword, err := s.database.InsertPassword(ctx, models.HashedPassword{
		Hash:       hashedPassword,
		Salt:       salt,
		Iterations: 600000,
	})
	if err != nil {
		return err
	}

	_, err = s.database.InsertKeys(ctx, models.UserKeyPair{
		UserID:     session.UserID,
		PasswordID: dbPassword.ID,
		Type:       models.KeyTypeRsa,
		Public:     x509.MarshalPKCS1PublicKey(&keyPair.Public),
		Private:    ciphertext,
	})
	session.PrivateKey = &keyPair.Private
	s.sessions[session.SessionID] = session

	return err
}

func (s *service) UnsealWithPassword(ctx context.Context, session *Session, password string) error {
	if session.PrivateKey != nil {
		return nil
	}

	keyPair, err := s.database.SelectKeys(ctx, session.UserID)
	if err != nil {
		return err
	}
	if keyPair.PasswordID == nil {
		return errors.New("no password associated with keys in database")
	}

	dbPassword, err := s.database.SelectPassword(ctx, *keyPair.PasswordID)
	if err != nil {
		return err
	}

	hashedPassword := hashPassword(password, dbPassword.Salt)

	if subtle.ConstantTimeCompare(dbPassword.Hash, hashedPassword) == 0 {
		return errors.New("password hash mismatch")
	}

	key := generateKeyFromPassword(password)
	cipher, err := aesCipher(key)
	if err != nil {
		return err
	}

	header := keyPair.Private[:8]
	if !slices.Equal(header, slices.Concat(MarkerBytes, AesGcmMarkerBytes)) {
		return errors.New("invalid encrypted private key")
	}
	nonceL := bytesToInt(keyPair.Private[8:12])
	nonce := keyPair.Private[12 : 12+nonceL]
	encryptedL := bytesToInt(keyPair.Private[12+nonceL : 12+nonceL+4])
	encrypted := keyPair.Private[12+nonceL+4 : 12+nonceL+4+encryptedL]
	decryptedPrivateKey, err := decrypt(encrypted, cipher, nonce)
	if err != nil {
		return err
	}

	parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(decryptedPrivateKey)
	if err != nil {
		return err
	}
	session.PrivateKey = parsedPrivateKey
	session.PrivateKey.Precompute()
	s.sessions[session.SessionID] = *session

	return nil
}

func hashPassword(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 600000, 32, sha256.New)
}

func generateKeyFromPassword(password string) []byte {
	return pbkdf2.Key([]byte(password), KeyDerivationSalt, 600000, 32, sha256.New)
}
