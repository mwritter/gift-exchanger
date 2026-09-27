-- +goose Up
ALTER TABLE exchanges
    ADD COLUMN exchange_description TEXT NOT NULL DEFAULT "",
    ADD COLUMN exchange_budget_cents INTEGER,
    ADD COLUMN exchange_invites TEXT[];


-- +goose Down
ALTER TABLE exchanges
    DROP COLUMN exchange_description,
    DROP COLUMN exchange_budget_cents,
    DROP COLUMN exchange_invites;