-- +goose Up
    -- unique(Timestamp, exchange_id, Base, Quote, Interval)
-- CREATE INDEX ON candle (timestamp desc);
-- CREATE INDEX ON candle (timestamp desc, id);

CREATE EXTENSION IF NOT EXISTS timescaledb;
SELECT create_hypertable('candle','timestamp');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE candle;

