-- +goose Up
-- +goose StatementBegin
CREATE TABLE threshold_secret_shares
(
    secret_id      BIGINT    NOT NULL REFERENCES secrets (id) ON DELETE CASCADE,
    user_id        BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    encrypted_share BYTEA    NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (secret_id, user_id)
);

CREATE INDEX threshold_secret_shares_user_id_idx ON threshold_secret_shares (user_id);

CREATE TABLE unlock_requests
(
    id           BIGSERIAL PRIMARY KEY,
    secret_id    BIGINT    NOT NULL REFERENCES secrets (id) ON DELETE CASCADE,
    requested_by BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at   TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '15 minutes',
    completed_at TIMESTAMP,
    UNIQUE (secret_id, requested_by, completed_at)
);

CREATE TABLE unlock_contributions
(
    request_id BIGINT    NOT NULL REFERENCES unlock_requests (id) ON DELETE CASCADE,
    user_id    BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    share      BYTEA     NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (request_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE unlock_contributions;
DROP TABLE unlock_requests;
DROP TABLE threshold_secret_shares;
-- +goose StatementEnd
