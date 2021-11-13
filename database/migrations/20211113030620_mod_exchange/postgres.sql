-- +goose Up
ALTER TABLE exchange ADD COLUMN balance double precision;
ALTER TABLE exchange ADD COLUMN num_pairs int;
ALTER TABLE exchange ADD COLUMN hq_country text;
ALTER TABLE exchange ADD COLUMN account_username text;

-- +goose Down
ALTER TABLE exchange DROP COLUMN balance;
ALTER TABLE exchange DROP COLUMN num_pairs;
ALTER TABLE exchange DROP COLUMN hq_country;
ALTER TABLE exchange DROP COLUMN account_username;
