package http

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/domain"
	"github.com/torfstack/synod/backend/models"
)

// mockDomainService is a flexible test double for domain.Service.
// Individual test functions assign only the stubs they need.
type mockDomainService struct {
	// UserService
	doesUserExistFn    func(ctx context.Context, username string) (bool, error)
	insertUserFn       func(ctx context.Context, user models.User) (models.ExistingUser, error)
	getUserFromTokenFn func(ctx context.Context, token *oidc.IDToken) (models.ExistingUser, error)

	// SecretService
	getSecretsFn   func(ctx context.Context, userID int64, cipher *crypto.AsymmetricCipher) ([]models.Secret, error)
	upsertSecretFn func(ctx context.Context, secret models.Secret, userID int64, cipher *crypto.AsymmetricCipher) (models.EncryptedSecret, error)

	// SessionService
	createSessionFn func(ctx context.Context, userID int64) (domain.Session, error)
	getSessionFn    func(token string) (*domain.Session, error)
	deleteSessionFn func(token string) error

	// SetupService
	isUserSetupFn           func(ctx context.Context, session domain.Session) (bool, error)
	setupUserPlainFn        func(ctx context.Context, session domain.Session) error
	setupUserWithPasswordFn func(ctx context.Context, session domain.Session, password crypto.Password) error
	unsealWithPasswordFn    func(ctx context.Context, session *domain.Session, password crypto.Password) error
}

var _ domain.Service = (*mockDomainService)(nil)

func (m *mockDomainService) DoesUserExist(ctx context.Context, username string) (bool, error) {
	if m.doesUserExistFn != nil {
		return m.doesUserExistFn(ctx, username)
	}
	return false, nil
}

func (m *mockDomainService) InsertUser(ctx context.Context, user models.User) (models.ExistingUser, error) {
	if m.insertUserFn != nil {
		return m.insertUserFn(ctx, user)
	}
	return models.ExistingUser{}, nil
}

func (m *mockDomainService) GetUserFromToken(ctx context.Context, token *oidc.IDToken) (models.ExistingUser, error) {
	if m.getUserFromTokenFn != nil {
		return m.getUserFromTokenFn(ctx, token)
	}
	return models.ExistingUser{}, nil
}

func (m *mockDomainService) GetSecrets(
	ctx context.Context,
	userID int64,
	cipher *crypto.AsymmetricCipher,
) ([]models.Secret, error) {
	if m.getSecretsFn != nil {
		return m.getSecretsFn(ctx, userID, cipher)
	}
	return nil, nil
}

func (m *mockDomainService) UpsertSecret(
	ctx context.Context,
	secret models.Secret,
	userID int64,
	cipher *crypto.AsymmetricCipher,
) (models.EncryptedSecret, error) {
	if m.upsertSecretFn != nil {
		return m.upsertSecretFn(ctx, secret, userID, cipher)
	}
	return models.EncryptedSecret{}, nil
}

func (m *mockDomainService) SearchShareRecipients(context.Context, int64, string) ([]models.ShareRecipient, error) {
	return nil, nil
}
func (m *mockDomainService) ShareSecret(context.Context, int64, int64, string, *crypto.AsymmetricCipher) error {
	return nil
}
func (m *mockDomainService) RevokeSecretAccess(context.Context, int64, int64, string) error {
	return nil
}
func (m *mockDomainService) GetSecretRecipients(context.Context, int64, int64) ([]models.ShareRecipient, error) {
	return nil, nil
}
func (m *mockDomainService) CreateThresholdSecret(context.Context, models.ThresholdSecretInput, int64) (int64, error) {
	return 1, nil
}
func (m *mockDomainService) StartUnlock(context.Context, int64, int64, *crypto.AsymmetricCipher) (int64, error) {
	return 1, nil
}
func (m *mockDomainService) GetUnlockRequests(context.Context, int64) ([]models.UnlockRequest, error) {
	return nil, nil
}
func (m *mockDomainService) ContributeToUnlock(context.Context, int64, int64, *crypto.AsymmetricCipher) error {
	return nil
}

func (m *mockDomainService) GetUnlockResult(
	context.Context,
	int64,
	int64,
	*crypto.AsymmetricCipher,
) (models.UnlockResult, error) {
	return models.UnlockResult{}, nil
}

func (m *mockDomainService) CreateSession(ctx context.Context, userID int64) (domain.Session, error) {
	if m.createSessionFn != nil {
		return m.createSessionFn(ctx, userID)
	}
	return domain.Session{}, nil
}

func (m *mockDomainService) GetSession(token string) (*domain.Session, error) {
	if m.getSessionFn != nil {
		return m.getSessionFn(token)
	}
	return nil, domain.ErrSessionNotFound
}

func (m *mockDomainService) DeleteSession(token string) error {
	if m.deleteSessionFn != nil {
		return m.deleteSessionFn(token)
	}
	return nil
}

func (m *mockDomainService) IsUserSetup(ctx context.Context, session domain.Session) (bool, error) {
	if m.isUserSetupFn != nil {
		return m.isUserSetupFn(ctx, session)
	}
	return false, nil
}

func (m *mockDomainService) SetupUserPlain(ctx context.Context, session domain.Session) error {
	if m.setupUserPlainFn != nil {
		return m.setupUserPlainFn(ctx, session)
	}
	return nil
}

func (m *mockDomainService) SetupUserWithPassword(
	ctx context.Context,
	session domain.Session,
	password crypto.Password,
) error {
	if m.setupUserWithPasswordFn != nil {
		return m.setupUserWithPasswordFn(ctx, session, password)
	}
	return nil
}

func (m *mockDomainService) UnsealWithPassword(
	ctx context.Context,
	session *domain.Session,
	password crypto.Password,
) error {
	if m.unsealWithPasswordFn != nil {
		return m.unsealWithPasswordFn(ctx, session, password)
	}
	return nil
}
