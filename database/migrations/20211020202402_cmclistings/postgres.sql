-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS cmc_latest_listing
(
    id SERIAL PRIMARY KEY,

	last_updated TIMESTAMP NOT NULL,
	market_cap DOUBLE PRECISION NOT NULL,
	market_cap_dominance DOUBLE PRECISION NOT NULL,
	percent_change_1h DOUBLE PRECISION NOT NULL,
	percent_change_24h DOUBLE PRECISION NOT NULL,
	percent_change_7d DOUBLE PRECISION NOT NULL,
	percent_change_30d DOUBLE PRECISION NOT NULL,
	percent_change_60d DOUBLE PRECISION NOT NULL,
	percent_change_90d DOUBLE PRECISION NOT NULL,
	percent_change_volume_24h DOUBLE PRECISION NOT NULL,
	percent_change_volume_30d DOUBLE PRECISION NOT NULL,
	percent_change_volume_7d DOUBLE PRECISION NOT NULL,
	price DOUBLE PRECISION NOT NULL,
	total_market_cap DOUBLE PRECISION NOT NULL,
	volume_24h DOUBLE PRECISION NOT NULL,
	volume_24h_reported DOUBLE PRECISION NOT NULL,
	volume_30d DOUBLE PRECISION NOT NULL,
	volume_30d_reported DOUBLE PRECISION NOT NULL,
	volume_7d DOUBLE PRECISION NOT NULL,
	volume_7d_reported DOUBLE PRECISION NOT NULL,
    circulating_supply DOUBLE PRECISION not null,
    cmc_rank int not null,
    date_added TIMESTAMP not null,
    fully_diluted_market_cap DOUBLE PRECISION NOT NULL,
    market_cap_by_total_supply DOUBLE PRECISION not null,
    max_supply DOUBLE PRECISION not null,
    name text not null,
    num_market_pairs int not null,
    slug text not null,
    symbol text not null,
    tags text,
    total_supply DOUBLE PRECISION not null,
    volume_change_24h DOUBLE PRECISION NOT NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc')
);
-- CREATE UNIQUE INDEX unique_symbols ON gateiocoin (base,quote);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE cmc_latest_listing;

