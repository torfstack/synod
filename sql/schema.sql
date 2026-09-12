CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    subject    TEXT      NOT NULL UNIQUE,
    email      TEXT      NOT NULL DEFAULT '',
    full_name  TEXT      NOT NULL DEFAULT '',
    sharing_id UUID      NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE secrets
(
    id             BIGSERIAL PRIMARY KEY,
    value          BYTEA     NOT NULL,
    key            TEXT      NOT NULL,
    url            TEXT      NOT NULL,
    tags           TEXT      NOT NULL,
    user_id        BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    secret_sharing INTEGER,
    envelope       BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE keys
(
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    password_id  BIGINT REFERENCES passwords (id) ON DELETE CASCADE,
    type         INT    NOT NULL,
    key_material BYTEA  NOT NULL,
    public_key   BYTEA
);

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

CREATE TABLE threshold_secret_shares
(
    secret_id       BIGINT    NOT NULL REFERENCES secrets (id) ON DELETE CASCADE,
    user_id         BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    encrypted_share BYTEA     NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
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

CREATE TABLE passwords
(
    id         BIGSERIAL PRIMARY KEY,
    hash       BYTEA  NOT NULL,
    salt       BYTEA  NOT NULL,
    iterations BIGINT NOT NULL
);
