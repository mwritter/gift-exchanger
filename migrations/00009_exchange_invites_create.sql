-- +goose Up
-- Keyed by email so an invite can exist before the invitee has a user row.
CREATE TABLE exchange_invites (
    exchange_id UUID NOT NULL REFERENCES exchanges (id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    accepted_at TIMESTAMPTZ,
    declined_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT exchange_invites_single_response_check CHECK (
        accepted_at IS NULL OR declined_at IS NULL
    )
);

CREATE UNIQUE INDEX exchange_invites_exchange_email_idx ON exchange_invites (exchange_id, lower(email));
CREATE INDEX exchange_invites_email_lower_idx ON exchange_invites (lower(email));

-- +goose Down
DROP TABLE IF EXISTS exchange_invites;
