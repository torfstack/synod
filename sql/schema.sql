CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    subject    TEXT      NOT NULL UNIQUE,
    email      TEXT      NOT NULL DEFAULT '',
    full_name  TEXT      NOT NULL DEFAULT '',
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
    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE keys
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    password_id BIGINT REFERENCES passwords (id) ON DELETE CASCADE,
    type        INT    NOT NULL,
    public      BYTEA  NOT NULL,
    private     BYTEA  NOT NULL
);

CREATE TABLE passwords
(
    id         BIGSERIAL PRIMARY KEY,
    hash       BYTEA  NOT NULL,
    salt       BYTEA  NOT NULL,
    iterations BIGINT NOT NULL
);