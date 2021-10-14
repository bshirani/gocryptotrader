-- +goose Up
-- +goose StatementBegin
CREATE TABLE live_order (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
    status text NOT NULL,
    order_type text NOT NULL,
    exchange text NOT NULL,
    strategy_id TEXT NOT NULL,
    internal_id text NOT NULL,
    side text,
    client_order_id text,
    amount REAL,
    symbol text,
    price REAL,
    fee REAL,
    cost REAL,
    submitted_at  timestamp,
    cancelled_at  timestamp,
    filled_at  timestamp,
    asset_type int,
    created_at  timestamp not null default CURRENT_TIMESTAMP,
    updated_at                    timestamp             NOT NULL default CURRENT_TIMESTAMP
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE live_order;
-- +goose StatementEnd
