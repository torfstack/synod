-- name: SelectSecrets :many
SELECT *
FROM secrets
WHERE user_id = $1;

-- name: InsertSecret :one
INSERT INTO secrets (value, key, url, tags, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

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

-- name: DoesUserExist :one
SELECT EXISTS(SELECT 1 FROM users WHERE subject = $1);

-- name: InsertUser :one
INSERT INTO users (subject, email, full_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: InsertKeys :one
INSERT INTO keys (user_id, type, public, private)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: SelectPublicKeyForUser :one
SELECT public
FROM keys
WHERE user_id = $1;