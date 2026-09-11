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
