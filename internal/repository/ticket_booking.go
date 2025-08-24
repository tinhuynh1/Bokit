package repository

import (
	"booking-svc/internal/domain"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ticketBookingRepo struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewTicketBookingRepo(db *gorm.DB, cache *redis.Client) domain.TicketBookingRepository {
	return &ticketBookingRepo{db: db, cache: cache}
}

func (r *ticketBookingRepo) AcquireLock(ctx context.Context, bookingIds []int) error {
	return r.cache.SetNX(ctx, "lock", "true", 15*time.Minute).Err()
}
func (r *ticketBookingRepo) CreateWithTx(ctx context.Context, tx *gorm.DB, booking *domain.TicketBooking) error {
	return tx.Create(booking).Error
}
func (r *ticketBookingRepo) GetByEventID(ctx context.Context, eventID int) ([]domain.TicketBooking, error) {
	var bookings []domain.TicketBooking
	err := r.db.Where("event_id = ?", eventID).Find(&bookings).Error
	return bookings, err
}
func (r *ticketBookingRepo) GetByEmail(ctx context.Context, email string) ([]domain.TicketBooking, error) {
	var bookings []domain.TicketBooking
	err := r.db.Where("email = ?", email).Find(&bookings).Error
	return bookings, err
}
func (r *ticketBookingRepo) GetByStatus(ctx context.Context, status string) ([]domain.TicketBooking, error) {
	return nil, nil
}
func (r *ticketBookingRepo) UpdateStatusByIds(ctx context.Context, tx *gorm.DB, bookingIds []int, status string) error {
	return tx.Model(&domain.TicketBooking{}).
		Where("id IN (?)", bookingIds).
		Update("status", status).
		Update("updated_at", time.Now()).
		Error
}

func (r *ticketBookingRepo) UpdateStatusById(ctx context.Context, bookingID int, status string) error {
	return r.db.Model(&domain.TicketBooking{}).Where("id = ?", bookingID).Update("status", status).Error
}

func (r *ticketBookingRepo) UpdateBooking(ctx context.Context, booking *domain.TicketBooking) error {
	return r.db.Model(&domain.TicketBooking{}).Where("id = ?", booking.ID).Updates(booking).Error
}

func (r *ticketBookingRepo) GetExpiredBooking(ctx context.Context, tx *gorm.DB) ([]domain.TicketBooking, error) {
	var bookings []domain.TicketBooking

	err := tx.Model(&domain.TicketBooking{}).
		Where("status = ? AND created_at < ?", domain.BookingStatusPending, time.Now().Add(-10*time.Minute)).
		Debug().
		Find(&bookings).Error
	return bookings, err
}

func (r *ticketBookingRepo) AcquireBookingLock(ctx context.Context) error {
	return r.cache.SetNX(ctx, "lock:booking_cancel", "true", 60*time.Second).Err()
}
func (r *ticketBookingRepo) ReleaseBookingLock(ctx context.Context) error {
	return r.cache.Del(ctx, "lock:booking_cancel").Err()
}

func (r *ticketBookingRepo) GetBookingsByEventIds(ctx context.Context, eventIds []int) ([]domain.TicketBooking, error) {
	var bookings []domain.TicketBooking
	err := r.db.Model(&domain.TicketBooking{}).Where("event_id IN (?)", eventIds).Find(&bookings).Error
	return bookings, err
}

func (r *ticketBookingRepo) GetBookingById(ctx context.Context, bookingID int) (domain.TicketBooking, error) {
	var booking domain.TicketBooking
	err := r.db.Model(&domain.TicketBooking{}).Where("id = ?", bookingID).First(&booking).Error
	return booking, err
}
