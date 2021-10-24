-- +goose Up
ALTER TABLE currency_pair_strategy
ADD COLUMN weight DOUBLE PRECISION NOT NULL DEFAULT 0;
-- +goose Down
ALTER TABLE currency_pair_strategy
DROP COLUMN weight;

