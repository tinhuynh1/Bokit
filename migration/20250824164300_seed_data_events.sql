-- +goose Up
-- +goose StatementBegin
-- 1. Concert Rock Festival
INSERT INTO events
(id, "name", description, date_time, ticket_price, available_tickets, sold_tickets, created_at, updated_at, deleted_at)
VALUES(nextval('events_id_seq'::regclass), 'Rock Festival 2025', 'Đêm nhạc rock hoành tráng với các ban nhạc nổi tiếng Việt Nam và quốc tế', '2025-09-15 19:00:00', 800000, 2000, 0, now(), now(), NULL);

-- 2. Tech Conference
INSERT INTO events
(id, "name", description, date_time, ticket_price, available_tickets, sold_tickets, created_at, updated_at, deleted_at)
VALUES(nextval('events_id_seq'::regclass), 'Vietnam Tech Summit 2025', 'Hội nghị công nghệ lớn nhất Việt Nam với các chuyên gia hàng đầu thế giới', '2025-10-20 08:00:00', 1200000, 500, 0, now(), now(), NULL);

-- 3. Art Exhibition
INSERT INTO events
(id, "name", description, date_time, ticket_price, available_tickets, sold_tickets, created_at, updated_at, deleted_at)
VALUES(nextval('events_id_seq'::regclass), 'Contemporary Art Exhibition', 'Triển lãm nghệ thuật đương đại với các tác phẩm từ 20 quốc gia', '2025-09-25 10:00:00', 150000, 800, 0, now(), now(), NULL);

-- 4. Food Festival
INSERT INTO events
(id, "name", description, date_time, ticket_price, available_tickets, sold_tickets, created_at, updated_at, deleted_at)
VALUES(nextval('events_id_seq'::regclass), 'Saigon Food Festival', 'Lễ hội ẩm thực Sài Gòn với 100+ món ăn truyền thống và hiện đại', '2025-11-10 17:00:00', 300000, 1500, 0, now(), now(), NULL);

-- 5. Startup Workshop
INSERT INTO events
(id, "name", description, date_time, ticket_price, available_tickets, sold_tickets, created_at, updated_at, deleted_at)
VALUES(nextval('events_id_seq'::regclass), 'Startup Success Workshop', 'Workshop khởi nghiệp với các CEO thành công và nhà đầu tư', '2025-12-05 14:00:00', 500000, 300, 0, now(), now(), NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM events;
-- +goose StatementEnd
