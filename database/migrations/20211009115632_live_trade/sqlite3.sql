-- +goose Up
CREATE TABLE live_trade (
    id int PRIMARY KEY,
    action order_action NOT NULL,
    entry_time timestamp NOT NULL,
    entry_price numeric NOT NULL,
    stop_loss_price numeric NOT NULL,
    take_profit_price numeric,
    exit_time timestamp NOT NULL,
    exit_price numeric NOT NULL,
    profit_pts numeric NOT NULL,
    risk numeric NOT NULL,
    risk_pts numeric NOT NULL,
    exit_reason exit_reason NOT NULL,
    units numeric,
    unit_profit numeric,
    unit_risk numeric,
    "prediction" numeric,
    net_profit numeric NOT NULL,
    weight numeric not null default 1,
    quantity numeric not null default 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
DROP TABLE live_trade;
