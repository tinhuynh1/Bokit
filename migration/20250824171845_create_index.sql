-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_ticket_bookings_event_id ON ticket_bookings (event_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_ticket_bookings_event_id;
-- +goose StatementEnd
