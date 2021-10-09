-- +goose Up
CREATE TABLE live_trade (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entry_price REAL NOT NULL,
    stop_loss_price REAL NOT NULL,
    take_profit_price REAL,
    exit_price REAL
);
-- +goose Down
DROP TABLE live_trade;
