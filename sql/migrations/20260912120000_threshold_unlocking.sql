-- +goose Up
-- +goose StatementBegin
CREATE TABLE threshold_secret_shares
(
    secret_id      BIGINT    NOT NULL REFERENCES secrets (id) ON DELETE CASCADE,
    user_id        BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    encrypted_share BYTEA    NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (secret_id, user_id)
);

CREATE INDEX threshold_secret_shares_user_id_idx ON threshold_secret_shares (user_id);

CREATE TABLE unlock_requests
(
    id           BIGSERIAL PRIMARY KEY,
    secret_id    BIGINT    NOT NULL REFERENCES secrets (id) ON DELETE CASCADE,
    requested_by BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at   TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '15 minutes',
    completed_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX unlock_requests_one_active_per_secret_idx
    ON unlock_requests (secret_id) WHERE completed_at IS NULL;

CREATE TABLE unlock_contributions
(
    request_id      BIGINT    NOT NULL REFERENCES unlock_requests (id) ON DELETE CASCADE,
    user_id         BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    encrypted_share BYTEA    NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (request_id, user_id)
);

CREATE TABLE threshold_unlock_grants
(
    request_id         BIGINT    NOT NULL REFERENCES unlock_requests (id) ON DELETE CASCADE,
    user_id            BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    encrypted_data_key BYTEA    NOT NULL,
    expires_at         TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (request_id, user_id)
);

CREATE INDEX threshold_unlock_grants_active_idx
    ON threshold_unlock_grants (user_id, expires_at);

ALTER TABLE secrets ADD CONSTRAINT threshold_secret_envelope_check
    CHECK (secret_sharing IS NULL OR (secret_sharing >= 2 AND envelope));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE secrets DROP CONSTRAINT threshold_secret_envelope_check;
DROP TABLE threshold_unlock_grants;
DROP TABLE unlock_contributions;
DROP TABLE unlock_requests;
DROP TABLE threshold_secret_shares;
-- +goose StatementEnd
