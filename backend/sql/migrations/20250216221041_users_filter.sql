-- +goose Up
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN username TO subject;
ALTER TABLE users ADD COLUMN email TEXT NOT NULL default '';
ALTER TABLE users ADD COLUMN full_name TEXT NOT NULL default '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN email;
ALTER TABLE users DROP COLUMN full_name;
ALTER TABLE users RENAME COLUMN subject TO username;
-- +goose StatementEnd
