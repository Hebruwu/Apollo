-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    username VARCHAR(255) NOT NULL PRIMARY KEY,
	email VARCHAR(255) NOT NULL UNIQUE,
	password_hash BYTEA NOT NULL,
	salt BYTEA NOT NULL,
	jwt_version INTEGER NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
