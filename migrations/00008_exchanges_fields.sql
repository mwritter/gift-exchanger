-- +goose Up
ALTER TABLE exchanges
    ADD COLUMN exchange_description TEXT NOT NULL DEFAULT '',
    ADD COLUMN exchange_budget_cents INTEGER CHECK (exchange_budget_cents >= 0);


-- +goose Down
ALTER TABLE exchanges
    DROP COLUMN exchange_description,
    DROP COLUMN exchange_budget_cents;
