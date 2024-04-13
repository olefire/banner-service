-- +goose Up
-- +goose StatementBegin
create table role_endpoints
(
    role text primary key,
    resource text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table role_endpoints;
-- +goose StatementEnd