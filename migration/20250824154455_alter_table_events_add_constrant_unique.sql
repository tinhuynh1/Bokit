-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ADD CONSTRAINT unique_event_id UNIQUE (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP CONSTRAINT unique_event_id;
-- +goose StatementEnd
