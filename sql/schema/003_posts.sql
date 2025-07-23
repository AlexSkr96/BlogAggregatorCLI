-- +goose Up
create table posts (
    id UUID primary key,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    title varchar(255),
    url varchar(255),
    description varchar(255),
    published_at timestamp,
    feed_id uuid not null references feeds(id)
);

-- +goose Down
drop table posts;
