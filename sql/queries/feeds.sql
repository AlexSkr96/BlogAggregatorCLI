-- name: CreateFeed :one
insert into feeds (id, name, url, user_id)
values ($1, $2, $3, $4)
returning id;

-- name: GetFeeds :many
select feeds.name, url, users.name from feeds
join users on feeds.user_id = users.id;
