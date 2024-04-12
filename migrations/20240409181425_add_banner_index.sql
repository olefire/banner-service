-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin

create unique index concurrently tag_feature_unique_idx on banner_feature_tag using hash (tag_id, feature_id);

-- +goose StatementEnd

-- +goose NO TRANSACTION
-- +goose Down
-- +goose StatementBegin
drop index concurrently tag_feature_unique_idx;
-- +goose StatementEnd
