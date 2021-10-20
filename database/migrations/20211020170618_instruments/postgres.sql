-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS instrument
(
    id SERIAL PRIMARY KEY,
    base varchar(30) NOT NULL,
    quote varchar(30) NOT NULL,
    market_cap decimal,
    data_from TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc')
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE instrument;

