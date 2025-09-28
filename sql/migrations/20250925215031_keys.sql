-- +goose Up
-- +goose StatementBegin
CREATE TABLE passwords
(
    id         BIGSERIAL PRIMARY KEY,
    hash       BYTEA  NOT NULL,
    salt       BYTEA  NOT NULL,
    iterations BIGINT NOT NULL
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE passwords;
DROP TABLE keys;
-- +goose StatementEnd
