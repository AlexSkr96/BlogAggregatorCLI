-- name: CreateUser :one
insert into users (id, created_at, updated_at, name)
values (
    $1,
    $2,
    $3,
    $4
)
returning *;

-- name: GetUserByUsername :one
select * from users
where name = $1;

-- name: DeleteAllUsers :exec
truncate table users;

-- name: GetUsers :many
select * from users;
