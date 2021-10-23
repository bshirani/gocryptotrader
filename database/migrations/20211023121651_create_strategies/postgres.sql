-- +goose Up
CREATE TYPE capture_type AS ENUM ('trend', 'trend2day', 'trend3day');
CREATE TYPE order_side AS ENUM ('BUY', 'SELL');

CREATE TABLE strategies (
    id SERIAL PRIMARY KEY,
    side order_side not null,
    capture capture_type not null,
    timeframe_days integer not null,
    created_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at       TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    UNIQUE(side, capture, timeframe_days)
);

INSERT INTO strategies (side, capture, timeframe_days) VALUES
    ('BUY', 'trend', 1),
    ('SELL', 'trend', 1),
    ('BUY', 'trend2day', 2),
    ('SELL', 'trend2day', 2),
    ('BUY', 'trend3day', 3),
    ('SELL', 'trend3day', 3)
;
-- +goose Down
DROP TABLE strategies;
DROP TYPE capture_type;
DROP TYPE order_side;
