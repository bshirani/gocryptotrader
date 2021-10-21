-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS instrument
(
    id SERIAL PRIMARY KEY,
    -- base varchar(30) NOT NULL,
    -- quote varchar(30) NOT NULL,
    symbol text NOT NULL,
    cmc_id integer NOT NULL,
    name text NOT NULL,
    slug text NOT NULL,
    first_historical_data TIMESTAMP NOT NULL,
    last_historical_data TIMESTAMP NOT NULL,
    market_cap decimal,
    active boolean NOT NULL,
    status boolean NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc')
);

-- CREATE UNIQUE INDEX unique_symbols ON instrument (base,quote);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE instrument;

