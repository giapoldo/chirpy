-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetSingleRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: RevokeSingleRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = $2 , updated_at = $3 WHERE token = $1;
