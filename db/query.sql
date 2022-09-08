-- name: CreateUser :exec
INSERT INTO users(tg_id, chat_id, role, first_name, last_name, user_name)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateUserRole :exec
UPDATE users SET role=$1
where tg_id=$2;

-- name: ReadUser :one
SELECT * FROM users
WHERE tg_id=$1;

-- name: ReadUsers :many
SELECT * FROM users
WHERE tg_id NOTNULL;

-- name: CreateFile :exec
INSERT INTO files(user_tg_id, tg_id, url)
VALUES ($1, $2, $3);

-- name: CreateMessage :exec
INSERT INTO messages(user_tg_id, text)
VALUES ($1, $2);