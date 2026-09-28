-- +goose Up
-- The login email carries both a link (token_hash) and a short code
-- (code_hash) so the person can finish logging in from the tab they started
-- in. Using either one consumes the row. attempts caps wrong code guesses.
ALTER TABLE magic_link_tokens
    ADD COLUMN code_hash TEXT,
    ADD COLUMN attempts INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE magic_link_tokens
    DROP COLUMN code_hash,
    DROP COLUMN attempts;
