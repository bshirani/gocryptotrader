-- +goose Up
CREATE TABLE live_signal (
    id SERIAL PRIMARY KEY,
    signal_time timestamp with time zone NOT NULL,
    valid_until timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    prediction double precision NOT NULL,
    strategy_name text NOT NULL,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL
);

-- +goose Down
DROP TABLE live_signal;
