-- +goose Up
-- +goose StatementBegin

CREATE TABLE banner
(
    banner_id       BIGSERIAL PRIMARY KEY,
    active_version  INTEGER   NOT NULL DEFAULT 1,
    is_active       BOOLEAN   NOT NULL DEFAULT TRUE,
    must_be_deleted BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner;

-- +goose StatementEnd

