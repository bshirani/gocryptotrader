-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS gateiocoin
(
    id SERIAL PRIMARY KEY,
    base varchar(30) NOT NULL,
    quote varchar(30) NOT NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc')
);
CREATE UNIQUE INDEX unique_symbols ON gateiocoin (base,quote);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE gateiocoin;

