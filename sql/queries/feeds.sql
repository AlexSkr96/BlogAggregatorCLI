-- name: CreateFeed :one
insert into feeds (id, name, url, user_id)
values ($1, $2, $3, $4)
returning id;

-- name: GetFeeds :many
select feeds.name, url, users.name from feeds
join users on feeds.user_id = users.id;

-- name: GetFeedByURL :one
select * from feeds
where url = $1;

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = now(), updated_at = now()
where id = $1;

-- name: FetchNextFeed :one
select * from feeds
order by last_fetched_at nulls first;

-- name: DeleteFeed :exec
delete from feeds
where id = $1;
