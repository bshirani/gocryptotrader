-- +goose Up
CREATE TYPE order_type AS ENUM ('MARKET', 'LIMIT', 'STOP');
CREATE TABLE public.live_order (
    id SERIAL PRIMARY KEY,
    status order_status NOT NULL,
    order_type order_type NOT NULL,
    exchange text NOT NULL,
    strategy_name text NOT NULL,
    internal_id text NOT NULL,
    side public.order_side NOT NULL,
    client_order_id text DEFAULT ''::text NOT NULL,
    amount double precision DEFAULT 0 NOT NULL,
    symbol text NOT NULL,
    price double precision DEFAULT 0 NOT NULL,
    stop_loss_price double precision DEFAULT 0 NOT NULL,
    take_profit_price double precision DEFAULT 0 NOT NULL,
    fee double precision DEFAULT 0 NOT NULL,
    cost double precision DEFAULT 0 NOT NULL,
    filled_at timestamp without time zone NOT NULL,
    asset_type integer DEFAULT 0 NOT NULL,
    submitted_at timestamp without time zone NOT NULL,
    cancelled_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL
);

CREATE TABLE public.live_trade (
    id SERIAL PRIMARY KEY,
    side public.order_side NOT NULL,
    entry_order_id integer NOT NULL,
    entry_price double precision NOT NULL,
    entry_time timestamp with time zone NOT NULL,
    exit_time timestamp with time zone NOT NULL,
    stop_loss_price double precision NOT NULL,
    strategy_name text NOT NULL,
    status text NOT NULL,
    amount double precision DEFAULT 0 NOT NULL,
    pair text NOT NULL,
    exchange text NOT NULL,
    take_profit_price double precision NOT NULL,
    profit_loss_points double precision NOT NULL,
    exit_price double precision NOT NULL,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL
);

ALTER TABLE ONLY public.live_trade
    ADD CONSTRAINT fk_live_trade_live_order_entry_id FOREIGN KEY (entry_order_id) REFERENCES public.live_order(id);

-- +goose Down
DROP TABLE live_trade;
DROP TABLE live_order;
