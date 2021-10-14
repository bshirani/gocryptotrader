-- +goose Up
CREATE TABLE live_trade (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
    side text not null,
    entry_order_id text NOT NULL,
    entry_price REAL NOT NULL,
    entry_time timestamp NOT NULL,
    stop_loss_price REAL NOT NULL,
    strategy_id TEXT NOT NULL,
    status TEXT NOT NULL,
    pair text NOT NULL,

    exit_time timestamp ,
    take_profit_price REAL,
    profit_loss_points REAL,
    exit_price REAL,

    created_at  timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL default CURRENT_TIMESTAMP
);
-- +goose Down
DROP TABLE live_trade;
