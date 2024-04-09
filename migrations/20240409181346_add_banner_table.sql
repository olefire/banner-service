-- +goose Up
-- +goose StatementBegin

CREATE TABLE  banner
(
    banner_id      BIGINT PRIMARY KEY,
    active_version INTEGER   NOT NULL,
    is_active      BOOLEAN            DEFAULT TRUE,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (banner_id) REFERENCES banner_feature_tag (banner_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner;

-- +goose StatementEnd

