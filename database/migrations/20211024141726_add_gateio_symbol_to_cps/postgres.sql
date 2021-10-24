-- +goose Up
ALTER TABLE currency_pair ADD COLUMN gateio_symbol text NOT NULL default '';
-- +goose Down
ALTER TABLE currency_pair DROP COLUMN gateio_symbol;
