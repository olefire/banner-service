-- +goose Up
-- +goose StatementBegin
CREATE TABLE banner_version
(
    banner_id  BIGINT    NOT NULL,
    version    SERIAL    NOT NULL,
    content    JSONB,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (banner_id) REFERENCES banner (banner_id),
    PRIMARY KEY (banner_id, version)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner_version;
-- +goose StatementEnd
