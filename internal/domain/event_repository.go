package domain

import (
	"context"

	"gorm.io/gorm"
)

type EventRepository interface {
	ListEvents(ctx context.Context, limit int, offset int, from string, to string) ([]Event, int64, error)
	GetEventsByIds(ctx context.Context, ids []int) ([]Event, error)
	GetEventForBooking(ctx context.Context, tx *gorm.DB, id int) (*Event, error)
	CreateEvent(ctx context.Context, event *Event) error
	UpdateEventWithTx(ctx context.Context, tx *gorm.DB, event *Event) error
	DeleteEvent(ctx context.Context, id int) error
	UpdateEventsWithTx(ctx context.Context, tx *gorm.DB, events []Event) error
	//GetEventStats(ctx context.Context, from string, to string) ([]dto.EventStatsResponse, error)
}
