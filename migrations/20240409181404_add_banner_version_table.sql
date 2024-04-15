-- +goose Up
-- +goose StatementBegin
create table banner_version
(
    banner_id  bigint    not null default 1,
    version    integer   not null default 1,
    content    jsonb              default '{}',
    updated_at timestamp not null default current_timestamp,
    primary key (banner_id, version)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table banner_version;
-- +goose StatementEnd
