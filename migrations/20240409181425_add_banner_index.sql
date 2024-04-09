-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin

CREATE UNIQUE INDEX CONCURRENTLY tag_feature_unique_idx ON banner_feature_tag (tag_id, feature_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX CONCURRENTLY tag_feature_unique_idx;
-- +goose StatementEnd
