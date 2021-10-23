-- +goose Up
CREATE TABLE live_order (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
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
    filled_at  timestamp,
    asset_type int,
    submitted_at  timestamp,
    cancelled_at  timestamp,
    created_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc')
);
-- +goose Down
DROP TABLE live_order;

alter table live_trade alter column entry_order_id type uuid USING entry_order_id::uuid;
ALTER TABLE live_trade ADD CONSTRAINT live_trade_live_order_fk FOREIGN KEY (entry_order_id) REFERENCES live_order (id);

