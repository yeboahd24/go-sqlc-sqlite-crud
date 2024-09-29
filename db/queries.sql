-- name: CreateUser :exec
INSERT INTO users (name, email)
VALUES (?, ?);

-- name: GetUserByID :one
SELECT id, name, email 
FROM users 
WHERE id = ?;

-- name: GetAllUsers :many
SELECT id, name, email 
FROM users;

-- name: UpdateUserByID :exec
UPDATE users
SET name = ?
WHERE id = ?;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = ?;

