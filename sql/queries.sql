-- name: SelectSecrets :many
SELECT *
FROM secrets
WHERE user_id = $1 AND envelope = FALSE;

-- name: SelectAccessibleSecrets :many
SELECT s.*, sa.encrypted_data_key, s.user_id = $1 AS owned
FROM secrets s
JOIN secret_access sa ON sa.secret_id = s.id
WHERE sa.user_id = $1;

-- name: InsertSecret :one
INSERT INTO secrets (value, key, url, tags, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: InsertSecretAccess :exec
INSERT INTO secret_access (secret_id, user_id, encrypted_data_key, granted_by)
VALUES ($1, $2, $3, $4)
ON CONFLICT (secret_id, user_id) DO NOTHING;

-- name: DeleteSecretAccess :execrows
DELETE FROM secret_access
WHERE secret_id = $1 AND user_id = $2;

-- name: SelectSecretRecipients :many
SELECT u.*
FROM users u
JOIN secret_access sa ON sa.user_id = u.id
JOIN secrets s ON s.id = sa.secret_id
WHERE s.id = $1 AND s.user_id = $2 AND u.id <> s.user_id
ORDER BY u.full_name, u.id;

-- name: SelectSecretForOwner :one
SELECT s.*, sa.encrypted_data_key
FROM secrets s
LEFT JOIN secret_access sa ON sa.secret_id = s.id AND sa.user_id = s.user_id
WHERE s.id = $1 AND s.user_id = $2;

-- name: UpdateSecretEnvelope :exec
UPDATE secrets
SET value = $1, key = '', url = '', tags = '', envelope = TRUE, updated_at = NOW()
WHERE id = $2 AND user_id = $3;

-- name: UpdateSecret :one
UPDATE secrets
SET value = $1,
    key   = $2,
    url   = $3,
    tags  = $4
WHERE user_id = $5
  AND id = $6
RETURNING *;

-- name: SelectUserByName :one
SELECT *
FROM users
WHERE subject = $1;

-- name: SearchUsers :many
SELECT *
FROM users
WHERE users.id <> $1
  AND (LOWER(full_name) LIKE LOWER($2) || '%' OR LOWER(email) = LOWER($2))
  AND EXISTS (SELECT 1 FROM keys WHERE keys.user_id = users.id AND keys.public_key IS NOT NULL)
ORDER BY CASE WHEN LOWER(email) = LOWER($2) THEN 0 ELSE 1 END, full_name, users.id
LIMIT 10;

-- name: SelectUserBySharingID :one
SELECT * FROM users WHERE sharing_id = $1;

-- name: DoesUserExist :one
SELECT EXISTS(SELECT 1 FROM users WHERE subject = $1);

-- name: InsertUser :one
INSERT INTO users (subject, email, full_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: InsertKeys :one
INSERT INTO keys (user_id, password_id, type, key_material, public_key)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: HasKeys :one
SELECT EXISTS(SELECT 1 FROM keys WHERE user_id = $1);

-- name: SelectKeys :one
SELECT *
FROM keys
WHERE user_id = $1;

-- name: UpdatePublicKey :exec
UPDATE keys SET public_key = $1 WHERE user_id = $2;

-- name: SelectPublicKey :one
SELECT id, public_key FROM keys WHERE user_id = $1 AND public_key IS NOT NULL;

-- name: SelectKeyMaterial :one
SELECT key_material
FROM keys
WHERE user_id = $1;

-- name: InsertPassword :one
INSERT INTO passwords (hash, salt, iterations)
VALUES ($1, $2, $3)
RETURNING *;

-- name: SelectPassword :one
SELECT *
FROM passwords
WHERE id = $1;

-- name: InsertThresholdSecret :one
INSERT INTO secrets (value, key, url, tags, user_id, secret_sharing, envelope)
VALUES ($1, $2, $3, $4, $5, $6, TRUE)
RETURNING *;

-- name: InsertThresholdShare :exec
INSERT INTO threshold_secret_shares (secret_id, user_id, encrypted_share)
VALUES ($1, $2, $3);

-- name: SelectThresholdSecrets :many
SELECT s.id, s.key, s.url, s.tags, s.user_id, s.secret_sharing
FROM secrets s
WHERE s.secret_sharing IS NOT NULL
  AND (s.user_id = $1 OR EXISTS (
      SELECT 1 FROM threshold_secret_shares tss WHERE tss.secret_id = s.id AND tss.user_id = $1
  ));

-- name: SelectThresholdSecretForParticipant :one
SELECT s.id, s.value, s.key, s.url, s.tags, s.user_id, s.secret_sharing, tss.encrypted_share
FROM secrets s
JOIN threshold_secret_shares tss ON tss.secret_id = s.id AND tss.user_id = $2
WHERE s.id = $1 AND s.secret_sharing IS NOT NULL;

-- name: InsertUnlockRequest :one
INSERT INTO unlock_requests (secret_id, requested_by)
SELECT s.id, $2 FROM secrets s
WHERE s.id = $1 AND s.secret_sharing IS NOT NULL
  AND (s.user_id = $2 OR EXISTS (
      SELECT 1 FROM threshold_secret_shares tss WHERE tss.secret_id = s.id AND tss.user_id = $2
  ))
RETURNING *;

-- name: SelectPendingUnlockRequests :many
SELECT ur.id, ur.secret_id, ur.requested_by, ur.created_at, ur.expires_at,
       s.key, s.secret_sharing, u.full_name AS requester_name,
       EXISTS (SELECT 1 FROM unlock_contributions uc WHERE uc.request_id = ur.id AND uc.user_id = $1) AS contributed,
       (SELECT COUNT(*) FROM unlock_contributions uc WHERE uc.request_id = ur.id) AS contributions
FROM unlock_requests ur
JOIN secrets s ON s.id = ur.secret_id
JOIN users u ON u.id = ur.requested_by
JOIN threshold_secret_shares tss ON tss.secret_id = ur.secret_id AND tss.user_id = $1
WHERE ur.completed_at IS NULL AND ur.expires_at > NOW()
ORDER BY ur.created_at DESC;

-- name: SelectUnlockRequest :one
SELECT ur.*, s.value AS encrypted_payload, s.secret_sharing
FROM unlock_requests ur
JOIN secrets s ON s.id = ur.secret_id
WHERE ur.id = $1 AND ur.requested_by = $2 AND ur.expires_at > NOW();

-- name: InsertUnlockContribution :exec
INSERT INTO unlock_contributions (request_id, user_id, share)
SELECT ur.id, $2, $3 FROM unlock_requests ur
JOIN threshold_secret_shares tss ON tss.secret_id = ur.secret_id AND tss.user_id = $2
WHERE ur.id = $1 AND ur.completed_at IS NULL AND ur.expires_at > NOW()
ON CONFLICT (request_id, user_id) DO NOTHING;

-- name: SelectUnlockContributions :many
SELECT share FROM unlock_contributions WHERE request_id = $1 ORDER BY user_id;

-- name: CompleteUnlockRequest :exec
UPDATE unlock_requests SET completed_at = NOW() WHERE id = $1;
