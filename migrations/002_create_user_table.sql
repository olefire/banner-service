-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "user"(
    username TEXT PRIMARY KEY,
    hash_password TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user";
-- +goose StatementEnd
