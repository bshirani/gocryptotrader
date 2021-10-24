-- +goose Up
-- +goose StatementBegin
drop table currency_pair_strategy;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
create table currency_pair_strategy ();
-- +goose StatementEnd
