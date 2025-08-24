package repository

import (
	"booking-svc/internal/domain"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type eventRepo struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewEventRepo(db *gorm.DB, cache *redis.Client) domain.EventRepository {
	return &eventRepo{db: db, cache: cache}
}

func (r *eventRepo) ListEvents(ctx context.Context, limit int, offset int, from string, to string) ([]domain.Event, int64, error) {
	var events []domain.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Debug().
		Select("id, name, description, date_time, ticket_price, available_tickets, sold_tickets, created_at, updated_at").
		Where("deleted_at IS NULL").
		Order("id DESC").
		Limit(limit).
		Offset(offset)

	if from != "" {
		query = query.Where("date_time >= ?", from)
	}

	if to != "" {
		query = query.Where("date_time < ?", to)
	}
	query = query.Count(&total)
	err := query.Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *eventRepo) GetEventForBooking(ctx context.Context, tx *gorm.DB, id int) (*domain.Event, error) {
	var event domain.Event
	err := tx.WithContext(ctx).Model(&domain.Event{}).
		Debug().
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepo) CreateEvent(ctx context.Context, event *domain.Event) error {
	err := r.db.WithContext(ctx).
		Debug().
		Create(event).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepo) UpdateEventWithTx(ctx context.Context, tx *gorm.DB, event *domain.Event) error {
	err := tx.WithContext(ctx).
		Debug().
		Save(event).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepo) DeleteEvent(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Model(&domain.Event{}).
		Debug().
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (r *eventRepo) GetEventsByIds(ctx context.Context, ids []int) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.WithContext(ctx).Model(&domain.Event{}).
		Debug().
		Where("id IN (?)", ids).
		Where("deleted_at IS NULL").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&events).Error
	return events, err
}

func (r *eventRepo) UpdateEventsWithTx(ctx context.Context, tx *gorm.DB, events []domain.Event) error {
	err := tx.WithContext(ctx).
		Debug().
		Save(&events).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepo) GetEventStats(ctx context.Context, from string, to string) ([]domain.Event, error) {
	var events []domain.Event
	query := r.db.WithContext(ctx).Model(&domain.Event{}).
		Debug().
		Where("deleted_at IS NULL")
	if from != "" {
		query = query.Where("date_time >= ?", from)
	}
	if to != "" {
		query = query.Where("date_time < ?", to)
	}
	err := query.Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, err
}
