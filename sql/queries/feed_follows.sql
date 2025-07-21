-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    insert into feed_follows (id, feed_id, user_id)
    values ($1, $2, $3)
    returning *
)
select
    iff.*,
    users.name,
    feeds.name
from inserted_feed_follow iff
inner join users on users.id = iff.user_id
inner join feeds on feeds.id = iff.feed_id;

-- name: GetFeedFollowsForUser :many
select
    ff.*,
    users.name as user_name,
    feeds.name as feed_name
from feed_follows ff
inner join users on users.id = ff.user_id
inner join feeds on feeds.id = ff.feed_id
where ff.user_id = $1;

-- name: DeleteFeedFollow :exec
delete from feed_follows
where user_id = $1
and feed_id = $2;
