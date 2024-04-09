-- +goose Up
-- +goose StatementBegin

CREATE TABLE users
(
    username      TEXT PRIMARY KEY,
    hash_password TEXT NOT NULL,
    role          TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
