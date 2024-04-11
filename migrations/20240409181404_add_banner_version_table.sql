-- +goose Up
-- +goose StatementBegin
CREATE TABLE banner_version
(
    banner_id  BIGINT    NOT NULL DEFAULT 1,
    version    INTEGER   NOT NULL DEFAULT 1,
    content    JSONB              DEFAULT '{}',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (banner_id, version)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner_version;
-- +goose StatementEnd
