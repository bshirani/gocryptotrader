-- +goose Up
-- +goose StatementBegin
alter table live_order alter column strategy_id type int USING strategy_id::integer;
ALTER TABLE live_trade alter COLUMN strategy_id type int USING strategy_id::integer;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE live_order ALTER COLUMN strategy_id type text;
ALTER TABLE live_trade ALTER COLUMN strategy_id type text;
-- +goose StatementEnd
