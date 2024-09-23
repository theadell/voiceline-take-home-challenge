-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES (?, ?)
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at, updated_at
FROM users
WHERE email = ?;

-- name: CreateUserProvider :one
INSERT INTO user_providers (user_id, provider, provider_user_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = ?, updated_at = CURRENT_TIMESTAMP
WHERE email = ?;

-- name: GetUserAndProviderInfo :one
SELECT 
    u.id AS user_id, 
    u.email AS user_email, 
    uop.id AS oauth_id
FROM users u
LEFT JOIN user_providers uop 
    ON u.id = uop.user_id 
    AND uop.provider = ? 
    AND uop.provider_user_id = ?
WHERE u.email = ?
LIMIT 1;
