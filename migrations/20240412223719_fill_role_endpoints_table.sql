-- +goose Up
-- +goose StatementBegin
insert into role_endpoints (role, resource)
values
    ('admin', '*'),
    ('user', 'GET /user_banner');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate role_endpoints;
-- +goose StatementEnd