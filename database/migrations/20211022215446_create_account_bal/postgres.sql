-- +goose Up
CREATE TABLE account_log (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    usd_balance DOUBLE PRECISION NOT NULL,
    btc_balance DOUBLE PRECISION NOT NULL,
    xrp_balance DOUBLE PRECISION NOT NULL,
    open_trades int not null,
    created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc')
);
-- +goose Down
DROP TABLE account_log;
