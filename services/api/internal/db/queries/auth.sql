-- name: InvalidateMagicLinkTokens :exec
UPDATE magic_link_tokens
SET consumed_at = now()
WHERE lower(email) = sqlc.arg(email) AND consumed_at IS NULL;

-- name: InsertMagicLinkToken :exec
INSERT INTO magic_link_tokens (email, token_hash, code_hash, expires_at)
VALUES (sqlc.arg(email), sqlc.arg(token_hash), sqlc.arg(code_hash), sqlc.arg(expires_at));

-- name: ConsumeMagicLinkToken :one
UPDATE magic_link_tokens
SET consumed_at = now()
WHERE token_hash = sqlc.arg(token_hash)
  AND consumed_at IS NULL
  AND expires_at > now()
RETURNING email;

-- name: GetActiveLoginCodeForUpdate :one
SELECT id, code_hash, attempts
FROM magic_link_tokens
WHERE lower(email) = sqlc.arg(email)
  AND consumed_at IS NULL
  AND expires_at > now()
  AND code_hash IS NOT NULL
ORDER BY created_at DESC
LIMIT 1
FOR UPDATE;

-- name: RecordLoginCodeAttempt :exec
UPDATE magic_link_tokens
SET attempts = attempts + 1,
    consumed_at = CASE WHEN attempts + 1 >= sqlc.arg(max_attempts)::int THEN now() ELSE consumed_at END
WHERE id = sqlc.arg(id);

-- name: ConsumeMagicLinkTokenByID :exec
UPDATE magic_link_tokens SET consumed_at = now() WHERE id = sqlc.arg(id);

-- name: UpsertUserByEmail :one
INSERT INTO users (email)
VALUES (sqlc.arg(email))
ON CONFLICT ((lower(email))) DO UPDATE SET email = users.email
RETURNING id;

-- name: InsertSession :exec
INSERT INTO sessions (user_id, token_hash, expires_at)
VALUES (sqlc.arg(user_id), sqlc.arg(token_hash), sqlc.arg(expires_at));

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token_hash = sqlc.arg(token_hash);

-- name: GetUserBySession :one
SELECT u.id, u.email, u.display_name
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.token_hash = sqlc.arg(token_hash) AND s.expires_at > now();
