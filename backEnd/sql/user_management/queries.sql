-- name: CreateUser :exec
INSERT INTO "user" (
    id, 
    username, 
    email, 
    password_hash, 
    date_of_birth
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5 
);

-- name: GetUser :one
SELECT id, username, bio, profile_picture_asset_id
FROM "user"
WHERE id = $1
LIMIT 1;