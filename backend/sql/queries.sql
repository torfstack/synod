-- name: SelectSecrets :many
SELECT * FROM secrets
WHERE user_id = $1;

-- name: InsertSecret :exec
INSERT INTO secrets (value, key, url, user_id)
VALUES ($1, $2, $3, $4);

-- name: UpdateSecret :exec
UPDATE secrets SET value = $1, key = $2, url = $3
WHERE user_id = $4 AND id = $5;

-- name: SelectUserByName :one
SELECT * FROM users
WHERE username = $1;

-- name: DoesUserExist :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: InsertUser :exec
INSERT INTO users (username)
VALUES ($1);