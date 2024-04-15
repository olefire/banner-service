-- +goose Up
-- +goose StatementBegin

create table users
(
    username      text primary key,
    hash_password text not null,
    role          text not null default 'user'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
