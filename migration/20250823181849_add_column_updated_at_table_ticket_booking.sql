-- +goose Up
-- +goose StatementBegin
ALTER TABLE ticket_bookings ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ticket_bookings DROP COLUMN updated_at;
-- +goose StatementEnd
