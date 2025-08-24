-- +goose Up
-- +goose StatementBegin
ALTER TABLE ticket_bookings ADD COLUMN total_price DOUBLE PRECISION NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ticket_bookings DROP COLUMN total_price;
-- +goose StatementEnd
