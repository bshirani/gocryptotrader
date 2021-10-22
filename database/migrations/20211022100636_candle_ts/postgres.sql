-- +goose Up
CREATE EXTENSION IF NOT EXISTS timescaledb;
SELECT create_hypertable('candle','timestamp');

-- +goose Down
DROP TABLE candle;

