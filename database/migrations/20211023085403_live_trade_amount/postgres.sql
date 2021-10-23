-- +goose Up
-- +goose StatementBegin
ALTER TABLE live_trade
ADD COLUMN amount DOUBLE PRECISION NOT NULL DEFAULT 0;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE live_trade DROP COLUMN amount;
-- +goose StatementEnd
