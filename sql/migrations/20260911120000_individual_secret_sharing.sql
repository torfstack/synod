-- +goose Up
-- +goose StatementBegin
ALTER TABLE keys ADD COLUMN public_key BYTEA;
ALTER TABLE users ADD COLUMN sharing_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE;
ALTER TABLE secrets ADD COLUMN envelope BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE secret_access
(
    secret_id          BIGINT    NOT NULL REFERENCES secrets (id) ON DELETE CASCADE,
    user_id            BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    encrypted_data_key BYTEA     NOT NULL,
    granted_by         BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at         TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (secret_id, user_id)
);

CREATE INDEX secret_access_user_id_idx ON secret_access (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secret_access;
ALTER TABLE secrets DROP COLUMN envelope;
ALTER TABLE keys DROP COLUMN public_key;
ALTER TABLE users DROP COLUMN sharing_id;
-- +goose StatementEnd
