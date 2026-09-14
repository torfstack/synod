package domain

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/db"
	"github.com/torfstack/synod/backend/logging"
	"github.com/torfstack/synod/backend/models"
)

var _ SecretService = &service{}

var ErrUnlockAlreadyActive = errors.New("another user is already unlocking this secret")

func (s *service) GetSecrets(
	ctx context.Context,
	userID int64,
	cipher *crypto.AsymmetricCipher,
) ([]models.Secret, error) {
	if cipher == nil {
		return []models.Secret{}, errors.New("need cipher to decrypt secrets")
	}
	encryptedSecrets, err := s.database.SelectSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}

	secrets := make([]models.Secret, 0, len(encryptedSecrets))
	for _, encryptedSecret := range encryptedSecrets {
		decrypted, err := s.DecryptSecret(ctx, encryptedSecret, cipher)
		if err != nil {
			return nil, err
		}
		decrypted.Owned = true
		secrets = append(secrets, decrypted)
	}
	accessible, err := s.database.SelectAccessibleSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}
	unlockedThreshold, err := s.database.SelectUnlockedThresholdSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}
	accessible = append(accessible, unlockedThreshold...)
	for _, stored := range accessible {
		wrappedKey, err := base64.StdEncoding.DecodeString(string(stored.EncryptedDataKey))
		if err != nil {
			return nil, err
		}
		key, err := cipher.Decrypt(wrappedKey)
		if err != nil {
			return nil, err
		}
		payloadCipher, err := crypto.SymmetricCipherFromKey(key)
		clear(key)
		if err != nil {
			return nil, err
		}
		payload, err := base64.StdEncoding.DecodeString(string(stored.EncryptedPayload))
		if err != nil {
			return nil, err
		}
		plaintext, err := payloadCipher.Decrypt(payload)
		if err != nil {
			return nil, err
		}
		var secret models.Secret
		if err := json.Unmarshal(plaintext, &secret); err != nil {
			return nil, err
		}
		secret.ID = &stored.ID
		secret.Owned = stored.Owned
		secret.UnlockedUntil = stored.UnlockedUntil
		secret.Threshold = stored.Threshold
		secret.Role = stored.Role
		secrets = append(secrets, secret)
	}
	thresholdSecrets, err := s.database.SelectThresholdSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, stored := range thresholdSecrets {
		id := stored.ID
		secrets = append(
			secrets,
			models.Secret{
				ID:        &id,
				Key:       stored.Key,
				Url:       stored.Url,
				Tags:      stored.Tags,
				Owned:     stored.OwnerID == userID,
				Locked:    true,
				Threshold: stored.Threshold,
				Role:      stored.Role,
			},
		)
	}
	return secrets, nil
}

func (s *service) UpsertSecret(
	ctx context.Context,
	secret models.Secret,
	userID int64,
	cipher *crypto.AsymmetricCipher,
) (models.EncryptedSecret, error) {
	if cipher == nil {
		return models.EncryptedSecret{}, errors.New("need cipher to encrypt secrets")
	}
	var result models.EncryptedSecret
	err := s.database.WithTx(ctx, func(database db.Database) error {
		thresholdSecret := false
		key, err := crypto.NewSymmetricKey()
		if err != nil {
			return err
		}
		if secret.ID != nil {
			stored, err := database.SelectSecretForManager(ctx, *secret.ID, userID)
			if err != nil {
				return err
			}
			thresholdSecret = stored.Threshold > 0
			if stored.Envelope {
				wrapped, err := base64.StdEncoding.DecodeString(string(stored.EncryptedDataKey))
				if err != nil {
					return err
				}
				key, err = cipher.Decrypt(wrapped)
				if err != nil {
					return err
				}
			}
		}
		defer clear(key)
		payloadCipher, err := crypto.SymmetricCipherFromKey(key)
		if err != nil {
			return err
		}
		plaintext, err := json.Marshal(secret)
		if err != nil {
			return err
		}
		payload, err := payloadCipher.Encrypt(plaintext)
		if err != nil {
			return err
		}
		encrypted := models.EncryptedSecret{
			ID:    secret.ID,
			Value: base64.StdEncoding.EncodeToString(payload),
			Key:   "",
			Url:   "",
			Tags:  []string{},
		}
		result, err = database.UpsertSecret(ctx, encrypted, userID)
		if err != nil {
			return err
		}
		if result.ID == nil {
			return errors.New("secret id missing after insert")
		}
		if err := database.UpdateSecretEnvelope(ctx, *result.ID, userID, []byte(result.Value)); err != nil {
			return err
		}
		if thresholdSecret {
			return nil
		}
		wrappedKey, err := cipher.Encrypt(key)
		if err != nil {
			return err
		}
		return database.InsertSecretAccess(
			ctx,
			*result.ID,
			userID,
			userID,
			[]byte(base64.StdEncoding.EncodeToString(wrappedKey)),
		)
	})
	if err != nil {
		logging.Errorf(ctx, "could not upsert secret: %v", err)
		return models.EncryptedSecret{}, errors.New("could not upsert secret")
	}
	return result, nil
}

func (s *service) SearchShareRecipients(
	ctx context.Context,
	userID int64,
	search string,
) ([]models.ShareRecipient, error) {
	users, err := s.database.SearchUsers(ctx, userID, search)
	if err != nil {
		return nil, err
	}
	return shareRecipients(users), nil
}

func shareRecipients(users []models.ExistingUser) []models.ShareRecipient {
	recipients := make([]models.ShareRecipient, len(users))
	for i, user := range users {
		recipients[i] = models.ShareRecipient{
			ID:       user.ID,
			Subject:  user.SharingID,
			Email:    maskEmail(user.Email),
			FullName: user.FullName,
		}
	}
	return recipients
}

func maskEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 || parts[0] == "" {
		return ""
	}
	return parts[0][:1] + "***@" + parts[1]
}

func (s *service) ShareSecret(
	ctx context.Context,
	secretID, ownerID int64,
	recipientSharingID string,
	cipher *crypto.AsymmetricCipher,
) error {
	if cipher == nil {
		return errors.New("need cipher to share secret")
	}
	return s.database.WithTx(ctx, func(database db.Database) error {
		stored, err := database.SelectSecretForManager(ctx, secretID, ownerID)
		if err != nil {
			return err
		}
		key, err := s.dataKeyForSharing(ctx, database, stored, ownerID, cipher)
		if err != nil {
			return err
		}
		defer clear(key)
		recipient, err := database.SelectUserBySharingID(ctx, recipientSharingID)
		if err != nil {
			return err
		}
		publicKey, err := database.SelectPublicKey(ctx, recipient.ID)
		if err != nil {
			return err
		}
		recipientCipher, err := crypto.AsymmetricCipherFromPublicKeyBytes(publicKey)
		if err != nil {
			return err
		}
		wrapped, err := recipientCipher.Encrypt(key)
		if err != nil {
			return err
		}
		return database.InsertSecretAccess(
			ctx,
			secretID,
			recipient.ID,
			ownerID,
			[]byte(base64.StdEncoding.EncodeToString(wrapped)),
		)
	})
}

func (s *service) dataKeyForSharing(
	ctx context.Context,
	database db.Database,
	stored models.AccessibleSecret,
	ownerID int64,
	cipher *crypto.AsymmetricCipher,
) ([]byte, error) {
	if stored.Envelope {
		wrapped, err := base64.StdEncoding.DecodeString(string(stored.EncryptedDataKey))
		if err != nil {
			return nil, err
		}
		return cipher.Decrypt(wrapped)
	}
	plain, err := s.DecryptSecret(ctx, stored.Legacy, cipher)
	if err != nil {
		return nil, err
	}
	key, err := crypto.NewSymmetricKey()
	if err != nil {
		return nil, err
	}
	payloadCipher, err := crypto.SymmetricCipherFromKey(key)
	if err != nil {
		return nil, err
	}
	encoded, err := json.Marshal(plain)
	if err != nil {
		return nil, err
	}
	payload, err := payloadCipher.Encrypt(encoded)
	if err != nil {
		return nil, err
	}
	wrapped, err := cipher.Encrypt(key)
	if err != nil {
		return nil, err
	}
	if err := database.UpdateSecretEnvelope(
		ctx,
		stored.ID,
		ownerID,
		[]byte(base64.StdEncoding.EncodeToString(payload)),
	); err != nil {
		return nil, err
	}
	if err := database.InsertSecretAccess(
		ctx,
		stored.ID,
		ownerID,
		ownerID,
		[]byte(base64.StdEncoding.EncodeToString(wrapped)),
	); err != nil {
		return nil, err
	}
	return key, nil
}

func (s *service) RevokeSecretAccess(ctx context.Context, secretID, ownerID int64, recipientSharingID string) error {
	_, err := s.database.SelectSecretForManager(ctx, secretID, ownerID)
	if err != nil {
		return err
	}
	recipient, err := s.database.SelectUserBySharingID(ctx, recipientSharingID)
	if err != nil {
		return err
	}
	if recipient.ID == ownerID {
		return errors.New("owner access cannot be revoked")
	}
	_, err = s.database.DeleteSecretAccess(ctx, secretID, recipient.ID)
	return err
}

func (s *service) GetSecretRecipients(ctx context.Context, secretID, ownerID int64) ([]models.ShareRecipient, error) {
	users, err := s.database.SelectSecretRecipients(ctx, secretID, ownerID)
	if err != nil {
		return nil, err
	}
	return shareRecipients(users), nil
}

func (s *service) CreateThresholdSecret(
	ctx context.Context,
	input models.ThresholdSecretInput,
	userID int64,
) (int64, error) {
	participantsInput := input.Participants
	if len(participantsInput) == 0 {
		participantsInput = make([]models.ThresholdParticipantInput, len(input.SharingIDs))
		for i, sharingID := range input.SharingIDs {
			participantsInput[i] = models.ThresholdParticipantInput{
				SharingID: sharingID, Role: models.ThresholdRoleHolder,
			}
		}
	}
	if input.Threshold < 2 || input.Threshold > len(participantsInput)+1 || input.Secret.ID != nil {
		return 0, errors.New("invalid threshold secret")
	}
	for _, participant := range participantsInput {
		if participant.Role != models.ThresholdRoleHolder && participant.Role != models.ThresholdRoleMaintainer {
			return 0, errors.New("invalid threshold participant role")
		}
	}
	key, err := crypto.NewSymmetricKey()
	if err != nil {
		return 0, err
	}
	defer clear(key)
	payloadCipher, err := crypto.SymmetricCipherFromKey(key)
	if err != nil {
		return 0, err
	}
	plaintext, err := json.Marshal(input.Secret)
	if err != nil {
		return 0, err
	}
	payload, err := payloadCipher.Encrypt(plaintext)
	if err != nil {
		return 0, err
	}
	shares, err := crypto.SplitSecret(key, input.Threshold, len(participantsInput)+1)
	if err != nil {
		return 0, err
	}
	returnID := int64(0)
	err = s.database.WithTx(ctx, func(database db.Database) error {
		participants := make([]models.ParticipantKey, 0, len(participantsInput)+1)
		participants = append(participants, models.ParticipantKey{UserID: userID, Role: models.ThresholdRoleOwner})
		seen := map[int64]struct{}{userID: {}}
		for _, participantInput := range participantsInput {
			participant, err := database.SelectUserBySharingID(ctx, participantInput.SharingID)
			if err != nil {
				return err
			}
			if _, exists := seen[participant.ID]; exists {
				return errors.New("duplicate threshold participant")
			}
			seen[participant.ID] = struct{}{}
			participants = append(participants, models.ParticipantKey{
				UserID: participant.ID, Role: participantInput.Role,
			})
		}
		stored := models.EncryptedSecret{
			Value: base64.StdEncoding.EncodeToString(payload),
			Key:   input.Secret.Key,
			Url:   input.Secret.Url,
			Tags:  input.Secret.Tags,
		}
		secretID, err := database.InsertThresholdSecret(ctx, stored, userID, input.Threshold)
		if err != nil {
			return err
		}
		returnID = secretID
		for i, participant := range participants {
			publicKey, err := database.SelectPublicKey(ctx, participant.UserID)
			if err != nil {
				return err
			}
			participantCipher, err := crypto.AsymmetricCipherFromPublicKeyBytes(publicKey)
			if err != nil {
				return err
			}
			encryptedShare, err := participantCipher.Encrypt(shares[i])
			if err != nil {
				return err
			}
			if err := database.InsertThresholdShare(
				ctx, secretID, participant.UserID, encryptedShare, participant.Role,
			); err != nil {
				return err
			}
		}
		return nil
	})
	return returnID, err
}

func (s *service) GetThresholdParticipants(
	ctx context.Context,
	secretID, userID int64,
) ([]models.ThresholdParticipant, error) {
	return s.database.SelectThresholdParticipants(ctx, secretID, userID)
}

func (s *service) AddThresholdParticipant(
	ctx context.Context,
	secretID, userID int64,
	input models.ThresholdParticipantInput,
	cipher *crypto.AsymmetricCipher,
) error {
	if input.Role != models.ThresholdRoleHolder && input.Role != models.ThresholdRoleMaintainer {
		return errors.New("invalid threshold participant role")
	}
	return s.changeThresholdParticipants(ctx, secretID, userID, cipher, func(
		database db.Database, actor models.AccessibleSecret, participants []models.ParticipantKey,
	) ([]models.ParticipantKey, error) {
		if actor.Role == models.ThresholdRoleMaintainer && input.Role != models.ThresholdRoleHolder {
			return nil, errors.New("only the owner can grant maintainer access")
		}
		user, err := database.SelectUserBySharingID(ctx, input.SharingID)
		if err != nil {
			return nil, err
		}
		for _, participant := range participants {
			if participant.UserID == user.ID {
				return nil, errors.New("threshold participant already exists")
			}
		}
		publicKey, err := database.SelectPublicKey(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		return append(participants, models.ParticipantKey{
			UserID: user.ID, PublicKey: publicKey, Role: input.Role,
		}), nil
	})
}

func (s *service) RemoveThresholdParticipant(
	ctx context.Context,
	secretID, userID int64,
	sharingID string,
	cipher *crypto.AsymmetricCipher,
) error {
	return s.changeThresholdParticipants(ctx, secretID, userID, cipher, func(
		database db.Database, actor models.AccessibleSecret, participants []models.ParticipantKey,
	) ([]models.ParticipantKey, error) {
		user, err := database.SelectUserBySharingID(ctx, sharingID)
		if err != nil {
			return nil, err
		}
		kept := make([]models.ParticipantKey, 0, len(participants)-1)
		found := false
		for _, participant := range participants {
			if participant.UserID != user.ID {
				kept = append(kept, participant)
				continue
			}
			found = true
			if participant.Role == models.ThresholdRoleOwner ||
				(actor.Role == models.ThresholdRoleMaintainer && participant.Role != models.ThresholdRoleHolder) {
				return nil, errors.New("participant cannot be removed")
			}
		}
		if !found {
			return nil, errors.New("threshold participant not found")
		}
		return kept, nil
	})
}

func (s *service) SetThresholdParticipantRole(
	ctx context.Context,
	secretID, userID int64,
	input models.ThresholdParticipantInput,
	cipher *crypto.AsymmetricCipher,
) error {
	if input.Role != models.ThresholdRoleHolder && input.Role != models.ThresholdRoleMaintainer {
		return errors.New("invalid threshold participant role")
	}
	return s.changeThresholdParticipants(ctx, secretID, userID, cipher, func(
		database db.Database, actor models.AccessibleSecret, participants []models.ParticipantKey,
	) ([]models.ParticipantKey, error) {
		if actor.Role != models.ThresholdRoleOwner {
			return nil, errors.New("only the owner can change participant roles")
		}
		user, err := database.SelectUserBySharingID(ctx, input.SharingID)
		if err != nil {
			return nil, err
		}
		found := false
		for i := range participants {
			if participants[i].UserID == user.ID {
				if participants[i].Role == models.ThresholdRoleOwner {
					return nil, errors.New("owner role cannot be changed")
				}
				participants[i].Role = input.Role
				found = true
			}
		}
		if !found {
			return nil, errors.New("threshold participant not found")
		}
		return participants, nil
	})
}

func (s *service) changeThresholdParticipants(
	ctx context.Context,
	secretID, userID int64,
	cipher *crypto.AsymmetricCipher,
	change func(db.Database, models.AccessibleSecret, []models.ParticipantKey) ([]models.ParticipantKey, error),
) error {
	if cipher == nil {
		return errors.New("secret must be unlocked to manage access")
	}
	return s.database.WithTx(ctx, func(database db.Database) error {
		stored, err := database.SelectSecretForManager(ctx, secretID, userID)
		if err != nil {
			return err
		}
		if stored.Threshold == 0 ||
			(stored.Role != models.ThresholdRoleOwner && stored.Role != models.ThresholdRoleMaintainer) {
			return errors.New("threshold access cannot be managed")
		}
		participants, err := database.SelectThresholdParticipantKeys(ctx, secretID)
		if err != nil {
			return err
		}
		participants, err = change(database, stored, participants)
		if err != nil {
			return err
		}
		if len(participants) < stored.Threshold {
			return errors.New("participant count cannot be lower than the threshold")
		}
		return s.rotateThresholdShares(ctx, database, stored, userID, participants, cipher)
	})
}

func (s *service) rotateThresholdShares(
	ctx context.Context,
	database db.Database,
	stored models.AccessibleSecret,
	userID int64,
	participants []models.ParticipantKey,
	cipher *crypto.AsymmetricCipher,
) error {
	wrapped, err := base64.StdEncoding.DecodeString(string(stored.EncryptedDataKey))
	if err != nil {
		return err
	}
	oldKey, err := cipher.Decrypt(wrapped)
	if err != nil {
		return err
	}
	defer clear(oldKey)
	oldCipher, err := crypto.SymmetricCipherFromKey(oldKey)
	if err != nil {
		return err
	}
	payload, err := base64.StdEncoding.DecodeString(string(stored.EncryptedPayload))
	if err != nil {
		return err
	}
	plaintext, err := oldCipher.Decrypt(payload)
	if err != nil {
		return err
	}
	newKey, err := crypto.NewSymmetricKey()
	if err != nil {
		return err
	}
	defer clear(newKey)
	newCipher, err := crypto.SymmetricCipherFromKey(newKey)
	if err != nil {
		return err
	}
	newPayload, err := newCipher.Encrypt(plaintext)
	if err != nil {
		return err
	}
	shares, err := crypto.SplitSecret(newKey, stored.Threshold, len(participants))
	if err != nil {
		return err
	}
	if err := database.UpdateSecretEnvelope(
		ctx, stored.ID, userID, []byte(base64.StdEncoding.EncodeToString(newPayload)),
	); err != nil {
		return err
	}
	if err := database.DeleteUnlockRequestsForSecret(ctx, stored.ID); err != nil {
		return err
	}
	if err := database.DeleteThresholdShares(ctx, stored.ID); err != nil {
		return err
	}
	for i, participant := range participants {
		participantCipher, err := crypto.AsymmetricCipherFromPublicKeyBytes(participant.PublicKey)
		if err != nil {
			return err
		}
		encryptedShare, err := participantCipher.Encrypt(shares[i])
		if err != nil {
			return err
		}
		if err := database.InsertThresholdShare(
			ctx, stored.ID, participant.UserID, encryptedShare, participant.Role,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) StartUnlock(
	ctx context.Context,
	secretID, userID int64,
	cipher *crypto.AsymmetricCipher,
) (int64, error) {
	if cipher == nil {
		return 0, errors.New("need cipher to unlock secret")
	}
	request, err := s.database.InsertUnlockRequest(ctx, secretID, userID)
	if err != nil {
		return 0, err
	}
	if request.RequesterID != userID {
		return 0, ErrUnlockAlreadyActive
	}
	if err := s.ContributeToUnlock(ctx, request.ID, userID, cipher); err != nil {
		return 0, err
	}
	return request.ID, nil
}

func (s *service) GetUnlockRequests(ctx context.Context, userID int64) ([]models.UnlockRequest, error) {
	return s.database.SelectPendingUnlockRequests(ctx, userID)
}

func (s *service) ContributeToUnlock(
	ctx context.Context,
	requestID, userID int64,
	cipher *crypto.AsymmetricCipher,
) error {
	if cipher == nil {
		return errors.New("need cipher to contribute share")
	}
	requests, err := s.database.SelectPendingUnlockRequests(ctx, userID)
	if err != nil {
		return err
	}
	var secretID int64
	var requesterID int64
	for _, request := range requests {
		if request.ID == requestID {
			secretID = request.SecretID
			requesterID = request.RequesterID
			break
		}
	}
	if secretID == 0 {
		return errors.New("unlock request not available")
	}
	stored, err := s.database.SelectThresholdSecretForParticipant(ctx, secretID, userID)
	if err != nil {
		return err
	}
	share, err := cipher.Decrypt(stored.EncryptedShare)
	if err != nil {
		return err
	}
	defer clear(share)
	requesterPublicKey, err := s.database.SelectPublicKey(ctx, requesterID)
	if err != nil {
		return err
	}
	requesterCipher, err := crypto.AsymmetricCipherFromPublicKeyBytes(requesterPublicKey)
	if err != nil {
		return err
	}
	encryptedShare, err := requesterCipher.Encrypt(share)
	if err != nil {
		return err
	}
	return s.database.InsertUnlockContribution(ctx, requestID, userID, encryptedShare)
}

func (s *service) GetUnlockResult(
	ctx context.Context,
	requestID, userID int64,
	cipher *crypto.AsymmetricCipher,
) (models.UnlockResult, error) {
	if cipher == nil {
		return models.UnlockResult{}, errors.New("need cipher to unlock secret")
	}
	stored, err := s.database.SelectUnlockRequest(ctx, requestID, userID)
	if err != nil {
		return models.UnlockResult{}, err
	}
	shares, err := s.database.SelectUnlockContributions(ctx, requestID)
	if err != nil {
		return models.UnlockResult{}, err
	}
	result := models.UnlockResult{Threshold: stored.Threshold, Contributions: len(shares)}
	if len(shares) < stored.Threshold {
		return result, nil
	}
	decryptedShares := make([][]byte, stored.Threshold)
	for i, encryptedShare := range shares[:stored.Threshold] {
		decryptedShares[i], err = cipher.Decrypt(encryptedShare)
		if err != nil {
			return models.UnlockResult{}, err
		}
		defer clear(decryptedShares[i])
	}
	key, err := crypto.CombineShares(decryptedShares)
	if err != nil {
		return models.UnlockResult{}, err
	}
	defer clear(key)
	payloadCipher, err := crypto.SymmetricCipherFromKey(key)
	if err != nil {
		return models.UnlockResult{}, err
	}
	payload, err := base64.StdEncoding.DecodeString(string(stored.EncryptedPayload))
	if err != nil {
		return models.UnlockResult{}, err
	}
	plaintext, err := payloadCipher.Decrypt(payload)
	if err != nil {
		return models.UnlockResult{}, err
	}
	var secret models.Secret
	if err := json.Unmarshal(plaintext, &secret); err != nil {
		return models.UnlockResult{}, err
	}
	secret.ID = &stored.ID
	secret.Threshold = stored.Threshold
	secret.Owned = stored.OwnerID == userID
	secret.Role = stored.Role
	result.Ready = true
	result.Secret = &secret
	expiresAt := time.Now().UTC().Add(10 * time.Minute)
	result.Secret.UnlockedUntil = &expiresAt
	err = s.database.WithTx(ctx, func(database db.Database) error {
		participants, err := database.SelectThresholdParticipantKeys(ctx, stored.ID)
		if err != nil {
			return err
		}
		for _, participant := range participants {
			participantCipher, err := crypto.AsymmetricCipherFromPublicKeyBytes(participant.PublicKey)
			if err != nil {
				return err
			}
			wrappedKey, err := participantCipher.Encrypt(key)
			if err != nil {
				return err
			}
			encodedKey := []byte(base64.StdEncoding.EncodeToString(wrappedKey))
			if err := database.InsertThresholdUnlockGrant(
				ctx, requestID, participant.UserID, encodedKey, expiresAt,
			); err != nil {
				return err
			}
		}
		return database.CompleteUnlockRequest(ctx, requestID)
	})
	if err != nil {
		return models.UnlockResult{}, err
	}
	return result, nil
}
