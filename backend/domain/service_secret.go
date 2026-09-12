package domain

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/torfstack/synod/backend/crypto"
	"github.com/torfstack/synod/backend/db"
	"github.com/torfstack/synod/backend/logging"
	"github.com/torfstack/synod/backend/models"
)

var _ SecretService = &service{}

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
		key, err := crypto.NewSymmetricKey()
		if err != nil {
			return err
		}
		if secret.ID != nil {
			stored, err := database.SelectSecretForOwner(ctx, *secret.ID, userID)
			if err != nil {
				return err
			}
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
		wrappedKey, err := cipher.Encrypt(key)
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
		stored, err := database.SelectSecretForOwner(ctx, secretID, ownerID)
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
	_, err := s.database.SelectSecretForOwner(ctx, secretID, ownerID)
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
	if input.Threshold < 2 || input.Threshold > len(input.SharingIDs)+1 || input.Secret.ID != nil {
		return 0, errors.New("invalid threshold secret")
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
	shares, err := crypto.SplitSecret(key, input.Threshold, len(input.SharingIDs)+1)
	if err != nil {
		return 0, err
	}
	returnID := int64(0)
	err = s.database.WithTx(ctx, func(database db.Database) error {
		participants := make([]models.ExistingUser, 0, len(input.SharingIDs)+1)
		participants = append(participants, models.ExistingUser{ID: userID})
		seen := map[int64]struct{}{userID: {}}
		for _, sharingID := range input.SharingIDs {
			participant, err := database.SelectUserBySharingID(ctx, sharingID)
			if err != nil {
				return err
			}
			if _, exists := seen[participant.ID]; exists {
				return errors.New("duplicate threshold participant")
			}
			seen[participant.ID] = struct{}{}
			participants = append(participants, participant)
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
			publicKey, err := database.SelectPublicKey(ctx, participant.ID)
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
			if err := database.InsertThresholdShare(ctx, secretID, participant.ID, encryptedShare); err != nil {
				return err
			}
		}
		return nil
	})
	return returnID, err
}

func (s *service) StartUnlock(
	ctx context.Context,
	secretID, userID int64,
	cipher *crypto.AsymmetricCipher,
) (int64, error) {
	if cipher == nil {
		return 0, errors.New("need cipher to unlock secret")
	}
	requestID, err := s.database.InsertUnlockRequest(ctx, secretID, userID)
	if err != nil {
		return 0, err
	}
	if err := s.ContributeToUnlock(ctx, requestID, userID, cipher); err != nil {
		return 0, err
	}
	return requestID, nil
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
	result.Ready = true
	result.Secret = &secret
	if err := s.database.CompleteUnlockRequest(ctx, requestID); err != nil {
		return models.UnlockResult{}, err
	}
	return result, nil
}
