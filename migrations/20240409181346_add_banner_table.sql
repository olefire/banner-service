-- +goose Up
-- +goose StatementBegin

create table banner
(
    banner_id       bigserial primary key,
    active_version  integer   not null default 1,
    is_active       boolean   not null default true,
    must_be_deleted boolean   not null default false,
    created_at      timestamp not null default current_timestamp
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table banner;

-- +goose StatementEnd

