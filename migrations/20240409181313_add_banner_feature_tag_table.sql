-- +goose Up
-- +goose StatementBegin

create table banner_feature_tag
(
    banner_id  bigint  not null,
    tag_id     integer not null,
    feature_id integer not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table banner_feature_tag;
-- +goose StatementEnd
