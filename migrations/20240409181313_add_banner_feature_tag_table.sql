-- +goose Up
-- +goose StatementBegin

CREATE TABLE banner_feature_tag
(
    banner_id  BIGINT  NOT NULL,
    tag_id     INTEGER NOT NULL,
    feature_id INTEGER NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner_feature_tag;
-- +goose StatementEnd
