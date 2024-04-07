-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXIST user(
    username TEXT PRIMARY KEY,
    hash_password VARCHAR(32) NOT NULL UNIQUE,
    is_admin BOOLEAN DEFAULT FALSE
);


CREATE TABLE IF NOT EXIST banner_version(
    banner_id INTEGER NOT NULL,
    version SERIAL NOT NULL,
    content TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (banner_id) REFERENCES banner(banner_id) ON DELETE CASCADE,
    PRIMARY KEY (banner_id, version)
);

CREATE TABLE IF NOT EXIST banner(
    banner_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    feature_id INTEGER NOT NULL,
    active_version INTEGER NOT NULL 1,
    PRIMARY KEY (banner_id)
);

CREATE UNIQUE INDEX banner_idx ON banner (tag_id, feature_id);

SELECT bv.content
FROM banner_version bv
JOIN banner b ON bv.banner_id = b.banner_id
WHERE b.tag_id = $1 AND b.feature_id = $2
  AND bv.active_version = TRUE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
DROP TABLE banner;
-- +goose StatementEnd
