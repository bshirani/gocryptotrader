-- +goose Up
alter table currency add column kraken_symbol text;
alter table currency add column kraken_active bool not null default false;

update currency set kraken_symbol='XBT' where symbol = 'BTC';
update currency set kraken_symbol=symbol where kraken_symbol is null;

CREATE TABLE currency_pair_strategy (
    id SERIAL PRIMARY KEY,
    currency_pair_id int NOT NULL,
    strategy_id int NOT NULL,
    side order_side NOT NULL,
    active boolean not null default false,
    created_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
     CONSTRAINT fk_currency_pair_strategy_currency_pair
      FOREIGN KEY(currency_pair_id)
	  REFERENCES currency(id),
     CONSTRAINT fk_currency_pair_strategy_strategy
      FOREIGN KEY(strategy_id)
	  REFERENCES strategy(id),
    UNIQUE(currency_pair_id, strategy_id, side)
);

INSERT INTO currency_pair_strategy(currency_pair_id, strategy_id, side)
SELECT cp.id, s.id, t.side
FROM currency_pair cp
CROSS JOIN strategy s
CROSS JOIN (select 'SELL'::order_side as side)t
ON CONFLICT DO NOTHING;

INSERT INTO currency_pair_strategy(currency_pair_id, strategy_id, side)
SELECT cp.id, s.id, t.side
FROM currency_pair cp
CROSS JOIN strategy s
CROSS JOIN (select 'BUY'::order_side as side)t
ON CONFLICT DO NOTHING;


-- INSERT INTO currency_pair_strategy(currency_pair_id, strategy_id, active) VALUES
--     (1,136,true),
--     (20,136,true)
-- ;

-- +goose Down
DROP TABLE currency_pair_strategy;
alter table currency drop column kraken_active;
alter table currency drop column kraken_symbol;
