-- Exchange IDs from the URL are compared as text so a malformed ID reads as
-- "not found" instead of a uuid parse error.

-- name: ListExchangesForUser :many
SELECT
    e.id,
    e.organizer_id,
    e.name,
    e.exchange_description,
    COALESCE(to_char(e.exchange_date, 'YYYY-MM-DD'), '')::text AS exchange_date,
    e.exchange_budget_cents,
    e.state
FROM exchanges e
JOIN exchange_members mine
    ON mine.exchange_id = e.id AND mine.user_id = sqlc.arg(user_id)
ORDER BY e.created_at DESC;

-- name: GetExchangeForUser :one
SELECT
    e.id,
    e.organizer_id,
    e.name,
    e.exchange_description,
    COALESCE(to_char(e.exchange_date, 'YYYY-MM-DD'), '')::text AS exchange_date,
    e.exchange_budget_cents,
    e.state
FROM exchanges e
JOIN exchange_members mine
    ON mine.exchange_id = e.id AND mine.user_id = sqlc.arg(user_id)
WHERE e.id::text = sqlc.arg(exchange_id)::text;

-- name: GetExchangeOrganizerID :one
SELECT organizer_id FROM exchanges WHERE id::text = sqlc.arg(exchange_id)::text;

-- name: ListOpenInviteEmails :many
SELECT email
FROM exchange_invites
WHERE exchange_id::text = sqlc.arg(exchange_id)::text
    AND accepted_at IS NULL
    AND declined_at IS NULL
ORDER BY lower(email);

-- name: ListMembersForExchanges :many
SELECT m.exchange_id, u.id AS user_id, u.email, u.display_name
FROM exchange_members m
JOIN users u ON u.id = m.user_id
WHERE m.exchange_id::text = ANY(sqlc.arg(exchange_ids)::text[])
ORDER BY m.created_at, lower(u.email);

-- name: InsertExchange :one
INSERT INTO exchanges (
    organizer_id, name, exchange_date, state, invite_token,
    exchange_description, exchange_budget_cents
)
VALUES (
    sqlc.arg(organizer_id),
    sqlc.arg(name),
    (sqlc.arg(exchange_date)::text)::date,
    sqlc.arg(state),
    sqlc.arg(invite_token),
    sqlc.arg(description),
    sqlc.narg(budget_cents)
)
RETURNING id;

-- name: InsertExchangeMember :exec
INSERT INTO exchange_members (exchange_id, user_id)
VALUES (sqlc.arg(exchange_id), sqlc.arg(user_id));

-- name: InsertExchangeInvite :exec
INSERT INTO exchange_invites (exchange_id, email)
VALUES (sqlc.arg(exchange_id), sqlc.arg(email));

-- name: UpdateExchange :execrows
UPDATE exchanges
SET name = sqlc.arg(name),
    exchange_description = sqlc.arg(description),
    exchange_date = (sqlc.arg(exchange_date)::text)::date,
    exchange_budget_cents = sqlc.narg(budget_cents),
    updated_at = now()
WHERE id::text = sqlc.arg(id)::text AND organizer_id = sqlc.arg(organizer_id);

-- name: ListInviteStates :many
SELECT
    lower(email)::text AS email,
    (accepted_at IS NOT NULL)::bool AS accepted,
    (declined_at IS NOT NULL)::bool AS declined
FROM exchange_invites
WHERE exchange_id::text = sqlc.arg(exchange_id)::text;

-- name: ReopenDeclinedInvite :exec
UPDATE exchange_invites
SET declined_at = NULL
WHERE exchange_id::text = sqlc.arg(exchange_id)::text AND lower(email) = sqlc.arg(email);

-- name: DeleteOpenInvite :exec
DELETE FROM exchange_invites
WHERE exchange_id::text = sqlc.arg(exchange_id)::text
    AND lower(email) = sqlc.arg(email)
    AND accepted_at IS NULL
    AND declined_at IS NULL;

-- name: DeleteExchange :execrows
DELETE FROM exchanges
WHERE id::text = sqlc.arg(id)::text AND organizer_id = sqlc.arg(organizer_id);
