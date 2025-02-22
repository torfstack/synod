-- +goose Up
-- +goose StatementBegin
ALTER TABLE secrets ADD COLUMN tags TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE secrets DROP COLUMN tags;
-- +goose StatementEnd
