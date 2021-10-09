-- +goose Up
CREATE TABLE live_trade (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
    entry_price REAL NOT NULL,
    stop_loss_price REAL NOT NULL,
    strategy_id TEXT NOT NULL,
    status TEXT NOT NULL,
    take_profit_price REAL,
    exit_price REAL
);
-- +goose Down
DROP TABLE live_trade;
