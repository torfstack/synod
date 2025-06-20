-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE secrets
(
    id BIGSERIAL PRIMARY KEY,
    value BYTEA NOT NULL,
    key TEXT NOT NULL,
    url TEXT NOT NULL,
    user_id BIGINT NOT NULL
        REFERENCES users (id)
        ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets;
DROP TABLE users;
-- +goose StatementEnd
