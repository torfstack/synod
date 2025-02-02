-- +goose Up
-- +goose StatementBegin
ALTER TABLE secrets ADD COLUMN secret_sharing INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE secrets DROP COLUMN secret_sharing;
-- +goose StatementEnd
