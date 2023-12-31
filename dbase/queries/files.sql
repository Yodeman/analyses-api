-- name: CreateFile :one
INSERT INTO files (
    username,
    data
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetFile :one
SELECT * FROM files
WHERE username = $1
LIMIT 1;

-- name: UpdateFile :one
UPDATE files
SET data = $1
WHERE username = $2
RETURNING *;
