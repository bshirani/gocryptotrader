-- +goose Up
CREATE TABLE live_trade (
    id SERIAL PRIMARY KEY,
    side text not null,
    entry_order_id text NOT NULL,
    entry_price double precision NOT NULL,
    entry_time TIMESTAMPTZ NOT NULL,
    stop_loss_price double precision NOT NULL,
    strategy_id TEXT NOT NULL,
    status TEXT NOT NULL,
    pair text NOT NULL,

    exit_time TIMESTAMPTZ ,
    take_profit_price REAL,
    profit_loss_points REAL,
    exit_price double precision,

    created_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc')
);
-- +goose Down
DROP TABLE live_trade;
