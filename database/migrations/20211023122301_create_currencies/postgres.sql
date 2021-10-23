-- +goose Up
ALTER TABLE cmc_coin RENAME TO currency;

CREATE TABLE currency_pair (
    id SERIAL PRIMARY KEY,
    base_id int NOT NULL,
    quote_id int NOT NULL,
    kraken_symbol text NOT NULL,
    active bool NOT NULL default false,
    created_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
     CONSTRAINT fk_currency_pair_cmc_coin_base
      FOREIGN KEY(base_id)
	  REFERENCES currency(id),
     CONSTRAINT fk_currency_pair_cmc_coin_quote
      FOREIGN KEY(quote_id)
	  REFERENCES currency(id)
);

INSERT INTO currency_pair(base_id, quote_id, kraken_symbol, active) VALUES
    (1,136,'XBT_USDT', true),
    (20,136,'XRP_USDT', true)
;

-- +goose Down
DROP TABLE currency_pair;
ALTER TABLE currency RENAME TO cmc_coin;
