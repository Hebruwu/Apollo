-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    username VARCHAR(255) NOT NULL PRIMARY KEY,
	email VARCHAR(255) NOT NULL UNIQUE,
	password_hash VARCHAR(255) NOT NULL,
	salt VARCHAR(32) NOT NULL,
	jwt_version INTEGER NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
