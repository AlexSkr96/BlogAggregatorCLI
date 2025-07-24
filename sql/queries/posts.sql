-- name: CreatePost :one
insert into posts (id, title, url, description, published_at, feed_id)
values ($1, $2, $3, $4, $5, $6)
returning id;

-- name: GetPostsForUser :many
select posts.* from posts
join feeds on feeds.id = posts.feed_id
join feed_follows ff on ff.feed_id = feeds.id
where ff.user_id = $1
order by published_at desc
limit $2;
