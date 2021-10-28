-- +goose Up
CREATE TYPE internal_order_type AS ENUM ('ENTER', 'TAKE_PROFIT', 'STOP_LOSS');
CREATE TYPE order_type AS ENUM ('MARKET', 'LIMIT', 'STOP');
create type order_status as enum ('NEW', 'ACTIVE', 'FILLED', 'CANCELLED');
CREATE TABLE public.live_order (
    id SERIAL PRIMARY KEY,
    status order_status NOT NULL,
    order_type order_type NOT NULL,
    internal_type internal_order_type NOT NULL,
    exchange text NOT NULL,
    strategy_name text NOT NULL,
    internal_id text NOT NULL,
    side public.order_side NOT NULL,
    client_order_id text DEFAULT ''::text NOT NULL,
    amount double precision NOT NULL,
    symbol text NOT NULL,
    price double precision NOT NULL,
    stop_loss_price double precision DEFAULT 0 NOT NULL,
    take_profit_price double precision DEFAULT 0 NOT NULL,
    fee double precision DEFAULT 0 NOT NULL,
    cost double precision DEFAULT 0 NOT NULL,
    filled_at timestamp without time zone,
    asset_type integer DEFAULT 0 NOT NULL,
    active_at timestamp without time zone,
    cancelled_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    CONSTRAINT strategy_name CHECK (
        strategy_name != ''
    ),
    CONSTRAINT new_check CHECK (
        status != 'NEW' OR
        (status = 'NEW' AND filled_at IS NULL AND active_at IS NULL AND cancelled_at IS NULL)
    ),
    CONSTRAINT cancelled_at CHECK (
        (cancelled_at IS NULL AND (status != 'CANCELLED'))
         OR (cancelled_at IS NOT NULL AND status IN ('CANCELLED'))
    ),
    CONSTRAINT filled_at CHECK (
        (filled_at IS NULL AND (status IN ('NEW', 'ACTIVE', 'CANCELLED')))
         OR (filled_at IS NOT NULL AND (status IN ('FILLED')))
    ),
    CONSTRAINT active_at CHECK (
        (active_at IS NULL AND (status != 'ACTIVE'))
         OR (active_at IS NOT NULL AND status IN ('ACTIVE', 'FILLED'))
    )
);

--
-- CREATE OR REPLACE FUNCTION check_internal_type_change()
--   RETURNS TRIGGER AS
-- $BODY$
-- BEGIN
--   RAISE EXCEPTION '"internal_type" column cannot get updated';
-- END;
-- $BODY$ LANGUAGE PLPGSQL;
--
-- CREATE TRIGGER live_order_update_trigger
-- BEFORE UPDATE OF "internal_type" ON "live_order"
-- FOR EACH ROW
-- WHEN (NEW."internal_type" IS DISTINCT FROM OLD."internal_type")
-- EXECUTE PROCEDURE check_internal_type_change();


CREATE TYPE trade_status AS ENUM ('OPEN', 'CLOSED');
CREATE TYPE exit_type AS ENUM ('STOP', 'TAKE_PROFIT');

CREATE TABLE public.live_trade (
    id SERIAL PRIMARY KEY,
    status trade_status NOT NULL,
    side order_side NOT NULL,
    entry_order_id integer NOT NULL UNIQUE,
    exit_order_id integer UNIQUE,
    entry_price double precision NOT NULL,
    exit_type exit_type,
    exit_price double precision ,
    entry_time timestamp with time zone NOT NULL,
    exit_time timestamp with time zone ,
    stop_loss_price double precision NOT NULL,
    strategy_name text NOT NULL,
    amount double precision DEFAULT 0 NOT NULL,
    risked_points double precision NOT NULL,
    risked_quote double precision NOT NULL,
    pair text NOT NULL,
    exchange text NOT NULL,
    take_profit_price double precision,
    profit_loss_points double precision ,
    profit_loss_quote double precision,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    CONSTRAINT risk_check CHECK(
        (risked_points != 0 AND risked_quote != 0)
    ),
    CONSTRAINT exit_price_time_check CHECK(
        (exit_time IS NULL and exit_price IS NULL ) OR
        (exit_time IS NOT NULL and exit_price IS NOT NULL)
    ),
    CONSTRAINT exit_type_check CHECK(
        (exit_time IS NULL AND exit_type IS NULL) OR
        (exit_time IS NOT NULL AND exit_type IS NOT NULL)
    ),
    CONSTRAINT exit_order_id_check CHECK(
        (exit_time IS NULL AND exit_order_id IS NULL) OR
        (exit_time IS NOT NULL AND exit_order_id IS NOT NULL)
    ),
    CONSTRAINT profit_check CHECK(
        (profit_loss_quote IS NULL AND profit_loss_points IS NULL AND exit_time IS NULL) OR
        (profit_loss_quote IS NOT NULL AND profit_loss_points IS NOT NULL AND exit_time IS NOT NULL)
    ),
    CONSTRAINT strategy_name CHECK(
        strategy_name != ''
    ),
    CONSTRAINT fk_live_trade_live_order_entry_id
        FOREIGN KEY (entry_order_id)
        REFERENCES live_order(id),
    CONSTRAINT fk_live_trade_live_order_exit_id
        FOREIGN KEY (exit_order_id)
        REFERENCES live_order(id)
);

-- +goose Down
DROP TABLE live_trade;
DROP TABLE live_order;
DROP TYPE order_type;
DROP TYPE order_status;
DROP TYPE trade_status;
DROP TYPE exit_type;
DROP TYPE internal_order_type;
-- DROP TRIGGER live_order_update_trigger;
-- DROP FUNCTION check_internal_type_change;;
