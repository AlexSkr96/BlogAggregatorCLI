-- +goose Up
create table feeds (
    id UUID primary key,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name varchar(255),
    url varchar(255) unique,
    user_id UUID references users(id) on delete cascade
);

-- +goose Down
drop table feeds;
