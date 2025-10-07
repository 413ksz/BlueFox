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