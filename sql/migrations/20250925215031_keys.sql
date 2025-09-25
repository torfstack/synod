-- +goose Up
-- +goose StatementBegin
CREATE TABLE keys
(
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    type    INT    NOT NULL,
    public  BYTEA  NOT NULL,
    private BYTEA  NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE keys;
-- +goose StatementEnd
