-- +goose Up
-- +goose StatementBegin
ALTER TABLE secret ADD COLUMN secret_charing INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE secret DROP COLUMN secret_charing;
-- +goose StatementEnd
