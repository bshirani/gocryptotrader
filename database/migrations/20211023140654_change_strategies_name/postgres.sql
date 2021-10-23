-- +goose Up
-- +goose StatementBegin
ALTER TABLE strategies RENAME TO strategy;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE strategy RENAME TO strategies;
-- +goose StatementEnd
