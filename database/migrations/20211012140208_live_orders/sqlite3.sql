-- +goose Up
-- +goose StatementBegin
CREATE TABLE live_order (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,

    status text NOT NULL,
    order_type text NOT NULL,
    exchange text NOT NULL,

    side text,
    client_order_id text,
    amount REAL,
    symbol text,
    price REAL,
    fee REAL,
    cost REAL,
    created_at  timestamp,
    updated_at timestamp,
    submitted_at  timestamp,
    cancelled_at  timestamp,
    filled_at  timestamp,
    asset_type int
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
