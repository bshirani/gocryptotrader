-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS cmc_latest_listing
(
    id SERIAL PRIMARY KEY,
    name text not null,
    symbol text not null,
    slug text not null,
    cmc_rank int not null,
    num_market_pairs int not null,
    circulating_supply decimal not null,
    total_supply DOUBLE PRECISION not null,
    market_cap_by_total_supply DOUBLE PRECISION not null,
    max_supply DOUBLE PRECISION not null,
    last_updated TIMESTAMP not null,
    date_added TIMESTAMP not null,
    tags text,
    base varchar(30) NOT NULL,
    quote varchar(30) NOT NULL,
    latest_price DOUBLE PRECISION not null,
    created_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc')
);
-- CREATE UNIQUE INDEX unique_symbols ON gateiocoin (base,quote);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE cmc_latest_listing;

