-- name: SelectSecrets :many
SELECT * FROM secrets
WHERE user_id = $1;

-- name: InsertSecret :exec
INSERT INTO secrets (value, key, url, tags, user_id)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateSecret :exec
UPDATE secrets SET value = $1, key = $2, url = $3, tags = $4
WHERE user_id = $5 AND id = $6;

-- name: SelectUserByName :one
SELECT * FROM users
WHERE username = $1;

-- name: DoesUserExist :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: InsertUser :exec
INSERT INTO users (username)
VALUES ($1);