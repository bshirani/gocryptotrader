-- +goose Up
CREATE TABLE live_order (
    id SERIAL PRIMARY KEY,
    status text NOT NULL,
    order_type text NOT NULL,
    exchange text NOT NULL,
    strategy_id TEXT NOT NULL,
    internal_id text NOT NULL,
    side order_side NOT NULL,
    client_order_id text NOT NULL DEFAULT '',
    amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    symbol text NOT NULL,
    price DOUBLE PRECISION NOT NULL DEFAULT 0,
    fee DOUBLE PRECISION NOT NULL DEFAULT 0,
    cost DOUBLE PRECISION NOT NULL DEFAULT 0REAL,
    filled_at timestamp NOT NULL,
    asset_type int NOT NULL DEFAULT 0,
    submitted_at  timestamp NOT NULL,
    cancelled_at  timestamp,
    created_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc')
);
-- +goose Down
DROP TABLE live_order;

