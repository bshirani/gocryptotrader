-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS candle
(
    ID SERIAL,
    exchange_id uuid REFERENCES exchange(id) NOT NULL,
    Base varchar(30) NOT NULL,
    Quote varchar(30) NOT NULL,
    Interval varchar(30) NOT NULL,
    Timestamp TIMESTAMPTZ NOT NULL,
    Open DOUBLE PRECISION NOT NULL,
    High DOUBLE PRECISION NOT NULL,
    Low DOUBLE PRECISION NOT NULL,
    Close DOUBLE PRECISION NOT NULL,
    Volume DOUBLE PRECISION NOT NULL,
    PRIMARY KEY (Timestamp, ID),
    unique(Timestamp, exchange_id, Base, Quote, Interval)
);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE candle;

