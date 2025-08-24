-- +goose Up
-- +goose StatementBegin
-- Xóa ràng buộc unique cũ
ALTER TABLE events DROP CONSTRAINT IF EXISTS unique_event_id;

-- Tạo partial unique index chỉ áp dụng cho records chưa bị xóa
CREATE UNIQUE INDEX idx_events_id_active ON events (id) 
WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events ADD CONSTRAINT unique_event_id UNIQUE (id);
DROP INDEX IF EXISTS idx_events_id_active;
-- +goose StatementEnd
