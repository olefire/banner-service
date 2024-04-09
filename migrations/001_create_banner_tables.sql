-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS banner_feature_tag(
     banner_id BIGINT PRIMARY KEY,
     tag_id INTEGER NOT NULL,
     feature_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS banner(
     banner_id BIGINT PRIMARY KEY,
     active_version INTEGER NOT NULL,
     is_active BOOLEAN DEFAULT TRUE,
     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
     FOREIGN KEY (banner_id) REFERENCES banner_feature_tag(banner_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS banner_version(
    banner_id BIGINT NOT NULL,
    version SERIAL NOT NULL,
    content JSONB,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (banner_id) REFERENCES banner(banner_id) ON DELETE CASCADE,
    PRIMARY KEY (banner_id, version)
);

CREATE UNIQUE INDEX banner_idx ON banner_feature_tag (tag_id, feature_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner_version;
DROP TABLE banner;
DROP TABLE banner_feature_tag;
DROP INDEX banner_idx;
-- +goose StatementEnd
