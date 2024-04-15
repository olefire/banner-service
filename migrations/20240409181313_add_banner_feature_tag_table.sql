-- +goose Up
-- +goose StatementBegin

create table banner_feature_tag
(
    banner_id  bigint  not null,
    tag_id     integer,
    feature_id integer,
    primary key (tag_id, feature_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table banner_feature_tag;
-- +goose StatementEnd
