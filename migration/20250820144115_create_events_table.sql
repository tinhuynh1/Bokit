-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE events_id_seq;
CREATE TABLE events (
  id             INTEGER NOT NULL DEFAULT nextval('events_id_seq'::regclass),
  name           TEXT NOT NULL,
  description    TEXT,
  date_time      TIMESTAMPTZ NOT NULL,
  ticket_price   DOUBLE PRECISION NOT NULL DEFAULT 0,
  total_tickets  INTEGER NOT NULL DEFAULT 0,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at     TIMESTAMPTZ DEFAULT NULL
);


CREATE SEQUENCE ticket_booking_id_seq;
CREATE TABLE ticket_bookings (
  id             INTEGER NOT NULL DEFAULT nextval('ticket_booking_id_seq'::regclass),
  event_id       INTEGER NOT NULL,
  email          TEXT NOT NULL,
  quantity       INTEGER NOT NULL,
  status         TEXT NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
DROP TABLE ticket_bookings;
DROP SEQUENCE events_id_seq;
DROP SEQUENCE ticket_booking_id_seq;
-- +goose StatementEnd
