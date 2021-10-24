-- +goose Up
ALTER TABLE strategy DROP COLUMN side;
ALTER TABLE strategy DROP COLUMN timeframe_days;
-- +goose Down
ALTER TABLE strategy ADD COLUMN side order_side;
ALTER TABLE strategy ADD COLUMN timeframe_days int;
