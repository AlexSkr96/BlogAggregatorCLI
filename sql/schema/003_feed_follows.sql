-- +goose Up
create table feed_follows (
    id UUID primary key,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    feed_id UUID NOT NULL REFERENCES feeds(id) on delete cascade,
    user_id UUID NOT NULL REFERENCES users(id) on delete cascade,
    unique (feed_id, user_id)
);

-- +goose Down
drop table feed_follows;
